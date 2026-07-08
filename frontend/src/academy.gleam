import academy/browser
import academy/icons
import gleam/dict
import gleam/dynamic/decode
import gleam/http/response
import gleam/int
import gleam/list
import gleam/option.{type Option, None, Some}
import gleam/result
import gleam/string
import gleam/time/calendar
import gleam/time/timestamp
import gleam/uri
import lustre
import lustre/attribute.{attribute, class}
import lustre/effect.{type Effect}
import lustre/element.{type Element}
import lustre/element/html
import lustre/element/keyed
import lustre/element/svg
import lustre/event
import modem
import rsvp

const base_url = ""

pub fn main() {
  let app = lustre.application(init, update, view)
  let assert Ok(_) = lustre.start(app, "#app", Nil)

  Nil
}

type Route {
  Threads
  Thread(id: String, content: option.Option(ModmailThread))
  Cases
  Case(id: Int, content: option.Option(AthenaCase))
  Stats
  Issues
  InterviewQuestions
  NotFound(path: List(String))
}

type WaveState {
  WaveInterviews
  WaveHelper
  WaveHistoric
}

type Trainee {
  Trainee(
    id: String,
    username: String,
    display_name: String,
    thread_participation_count: Int,
    message_count: Int,
    case_count: Int,
  )
}

type Issue {
  Issue(
    id: Int,
    created_by: Option(Int),
    trainee_id: Option(Int),
    thread_id: Option(Int),
    message_id: Option(Int),
    created_at: timestamp.Timestamp,
    status: String,
    reason: String,
    category: String,
  )
}

fn issue_decoder() -> decode.Decoder(Issue) {
  use id <- decode.field("id", decode.int)
  use created_by <- decode.field("created_by", decode.optional(decode.int))
  use trainee_id <- decode.field("trainee_id", decode.optional(decode.int))
  use thread_id <- decode.field("thread_id", decode.optional(decode.int))
  use message_id <- decode.field("message_id", decode.optional(decode.int))
  use created_at <- decode.field("created_at", timestamp_decoder())
  use status <- decode.field("status", decode.string)
  use reason <- decode.field("reason", decode.string)
  use category <- decode.field("category", decode.string)

  decode.success(Issue(
    id:,
    created_by:,
    trainee_id:,
    thread_id:,
    message_id:,
    created_at:,
    status:,
    reason:,
    category:,
  ))
}

fn trainee_decoder() -> decode.Decoder(Trainee) {
  use id <- decode.field("snowflake", decode.string)
  use username <- decode.field("username", decode.string)
  use display_name <- decode.field("display_name", decode.string)
  // use thread_participation_count <- decode.field(
  //   "thread_participation_count",
  //   decode.int,
  // )
  // use message_count <- decode.field("message_count", decode.int)
  // use case_count <- decode.field("case_count", decode.int)

  decode.success(Trainee(
    id:,
    username:,
    display_name:,
    thread_participation_count: 0,
    message_count: 0,
    case_count: 0,
  ))
}

fn wave_state_decoder() -> decode.Decoder(WaveState) {
  use variant <- decode.then(decode.string)
  case variant {
    "interviews" -> decode.success(WaveInterviews)
    "helper" -> decode.success(WaveHelper)
    "historic" -> decode.success(WaveHistoric)
    _ -> decode.failure(WaveInterviews, "WaveState")
  }
}

type Wave {
  Wave(
    id: Int,
    state: WaveState,
    created_at: timestamp.Timestamp,
    begin_at: timestamp.Timestamp,
    close_at: timestamp.Timestamp,
    trainees: List(Trainee),
  )
}

fn timestamp_decoder() -> decode.Decoder(timestamp.Timestamp) {
  use num <- decode.then(decode.int)
  decode.success(timestamp.from_unix_seconds(num))
}

fn wave_decoder() -> decode.Decoder(Wave) {
  use id <- decode.field("id", decode.int)
  use state <- decode.field("state", wave_state_decoder())
  use created_at <- decode.field("created_at", timestamp_decoder())
  use begin_at <- decode.field("begin_at", timestamp_decoder())
  use close_at <- decode.field("close_at", timestamp_decoder())
  use trainees <- decode.field("trainees", decode.list(trainee_decoder()))

  decode.success(Wave(id:, state:, created_at:, begin_at:, close_at:, trainees:))
}

type ThreadMsgKind {
  InternalMsg
  IncomingMsg
  OutgoingMsg
  SystemMsg
  CommandMsg
}

type ThreadRole {
  TraineeRole
  ModRole
  HelperRole
  AdminRole
  SystemRole
  UserRole
}

type ThreadMsg {
  ThreadMsg(
    id: Int,
    kind: ThreadMsgKind,
    role: ThreadRole,
    anonymous: Bool,
    user_id: String,
    user_name: String,
    created_at: Int,
    body: String,
    attachments: List(String),
  )
}

fn message_kind_decoder() -> decode.Decoder(ThreadMsgKind) {
  use num <- decode.then(decode.int)
  case num {
    1 -> decode.success(SystemMsg)
    2 -> decode.success(InternalMsg)
    3 -> decode.success(IncomingMsg)
    4 -> decode.success(OutgoingMsg)
    6 -> decode.success(CommandMsg)
    _ -> decode.success(SystemMsg)
    // decode.failure(SystemMsg, "MessageKind decoder")
  }
}

fn thread_msg_decoder() -> decode.Decoder(ThreadMsg) {
  use id <- decode.field("id", decode.int)
  use kind <- decode.field("kind", message_kind_decoder())
  use anonymous <- decode.field(
    "anonymous",
    decode.one_of(decode.map(decode.int, int.is_odd), [decode.bool]),
  )
  use role <- decode.field(
    "role",
    decode.then(decode.string, fn(str) {
      case string.lowercase(str) {
        "admin" -> AdminRole
        "helper" -> HelperRole
        "moderator" | "mod" -> ModRole
        _ -> SystemRole
      }
      |> decode.success
    }),
  )
  use user_id <- decode.field("user_id", decode.string)
  use user_name <- decode.field("user_name", decode.string)
  use created_at <- decode.field("created_at", decode.int)
  use body <- decode.field("body", decode.string)
  use attachments <- decode.field("attachments", decode.list(decode.string))

  decode.success(ThreadMsg(
    id:,
    kind:,
    role:,
    anonymous:,
    user_id:,
    user_name:,
    body:,
    created_at:,
    attachments:,
  ))
}

type ThreadStatus {
  ThreadOpen
  ThreadClosed
  ThreadSuspended
}

fn thread_status_decoder() -> decode.Decoder(ThreadStatus) {
  use variant <- decode.then(decode.int)
  case variant {
    1 -> decode.success(ThreadOpen)
    2 -> decode.success(ThreadClosed)
    3 -> decode.success(ThreadSuspended)
    _ -> decode.failure(ThreadOpen, "ThreadStatus")
  }
}

type ModmailThread {
  ModmailThread(
    id: String,
    user_name: String,
    user_id: String,
    status: ThreadStatus,
    inbound_messages: Int,
    outbound_messages: Int,
    chat_messages: Int,
    roles: List(String),
    participants: List(String),
    messages: List(ThreadMsg),
  )
}

fn modmail_thread_decoder() -> decode.Decoder(ModmailThread) {
  use id <- decode.field("id", decode.string)
  use user_name <- decode.field("user_name", decode.string)
  use user_id <- decode.field("user_id", decode.string)
  use status <- decode.field("status", thread_status_decoder())
  use inbound_messages <- decode.field("inbound_messages", decode.int)
  use outbound_messages <- decode.field("outbound_messages", decode.int)
  use chat_messages <- decode.field("chat_messages", decode.int)
  use participants <- decode.field("participants", decode.list(decode.string))
  use messages <- decode.optional_field(
    "messages",
    [],
    decode.list(thread_msg_decoder()),
  )
  use roles <- decode.field("roles", decode.list(decode.string))
  let roles = list.filter(roles, fn(r) { r != "@everyone" })

  decode.success(ModmailThread(
    id:,
    user_name:,
    user_id:,
    status:,
    inbound_messages:,
    outbound_messages:,
    chat_messages:,
    participants:,
    roles:,
    messages:,
  ))
}

type AthenaCaseNote {
  AthenaCaseNote(id: Int, mod_id: String, body: String)
}

fn athena_case_note_decoder() -> decode.Decoder(AthenaCaseNote) {
  use id <- decode.field("id", decode.int)
  use mod_id <- decode.field("mod_id", decode.string)
  use body <- decode.field("body", decode.string)
  decode.success(AthenaCaseNote(id:, mod_id:, body:))
}

type AthenaCaseKind {
  CaseBan
  CaseUnban
  CaseNote
  CaseWarn
  CaseKick
  CaseMute
  CaseUnmute
  CaseDeleted
  CaseSoftban
}

fn athena_case_kind_decoder() -> decode.Decoder(AthenaCaseKind) {
  use variant <- decode.then(decode.int)
  case variant {
    1 -> decode.success(CaseBan)
    2 -> decode.success(CaseUnban)
    3 -> decode.success(CaseNote)
    4 -> decode.success(CaseWarn)
    5 -> decode.success(CaseKick)
    6 -> decode.success(CaseMute)
    7 -> decode.success(CaseUnmute)
    8 -> decode.success(CaseDeleted)
    9 -> decode.success(CaseSoftban)
    _ -> decode.failure(CaseBan, "AthenaCaseKind")
  }
}

type AthenaCase {
  AthenaCase(
    id: Int,
    case_number: Int,
    mod_id: String,
    mod_name: String,
    actioned_user_id: String,
    actioned_user_name: String,
    created_at: timestamp.Timestamp,
    kind: AthenaCaseKind,
    notes: List(AthenaCaseNote),
  )
}

fn athena_case_decoder() -> decode.Decoder(AthenaCase) {
  use id <- decode.field("id", decode.int)
  use case_number <- decode.field("case_number", decode.int)
  use mod_id <- decode.field("mod_id", decode.string)
  use mod_name <- decode.field("mod_name", decode.string)
  use actioned_user_id <- decode.field("actioned_user_id", decode.string)
  use actioned_user_name <- decode.field("actioned_user_name", decode.string)
  use created_at <- decode.field("created_at", timestamp_decoder())
  use kind <- decode.field("type", athena_case_kind_decoder())
  use notes <- decode.optional_field(
    "notes",
    [],
    decode.list(athena_case_note_decoder()),
  )
  decode.success(AthenaCase(
    id:,
    case_number:,
    mod_id:,
    mod_name:,
    actioned_user_id:,
    actioned_user_name:,
    created_at:,
    kind:,
    notes:,
  ))
}

type UserRole {
  AdminUser
  HelperUser
  ModUser
  TraineeUser
  UnknownUser
}

fn user_role_decoder() -> decode.Decoder(UserRole) {
  use variant <- decode.then(decode.string)
  case variant {
    "admin" -> decode.success(AdminUser)
    "helper" -> decode.success(HelperUser)
    "mod" -> decode.success(ModUser)
    "trainee" -> decode.success(TraineeUser)
    _ -> decode.failure(UnknownUser, "UserRole")
  }
}

type User {
  User(id: String, avatar_hash: String, display_name: String, role: UserRole)
}

fn user_decoder() -> decode.Decoder(User) {
  use id <- decode.field("snowflake", decode.string)
  use display_name <- decode.field("display_name", decode.string)
  use avatar_hash <- decode.field("avatar_hash", decode.string)
  use role <- decode.field("role", user_role_decoder())
  decode.success(User(id:, avatar_hash:, display_name:, role:))
}

type Toast {
  ToastError(msg: String)
  ToastSuccess(msg: String)
  ToastWarning(msg: String)
}

type Modal {
  ClosedModal
  ThreadIssueModal(thread_id: String, message_id: Int, mod_id: String)
}

type Model {
  Model(
    route: Route,
    authenticated: Bool,
    toasts: dict.Dict(Int, Toast),
    wave: Wave,
    loading: Bool,
    user: User,
    issues: List(Issue),
    trainees: List(Trainee),
    threads: List(ModmailThread),
    total_cases: Int,
    total_threads: Int,
    total_issues: Int,
    cases: List(AthenaCase),
    // Interview stuff
    interview_questions: List(String),
    // Thread filtering
    thread_filter: String,
    threads_open: Bool,
    threads_closed: Bool,
    thread_trainee: Option(String),
    view_commands: Bool,
    searching: Bool,
    modal: Modal,
  )
}

fn init(_) -> #(Model, Effect(Message)) {
  let route =
    modem.initial_uri()
    |> result.map(fn(uri) { uri.path_segments(uri.path) })
    |> fn(path) {
      case path {
        Ok([]) | Ok(["threads"]) -> Threads

        Ok(["threads", id]) -> Thread(id, option.None)

        Ok(["cases"]) -> Cases

        Ok(["questions"]) -> InterviewQuestions

        Ok(["cases", id] as path) ->
          case int.parse(id) {
            Ok(id) -> Case(id, option.None)
            _ -> NotFound(path)
          }

        Ok(["stats"]) -> Stats

        Ok(["issues"]) -> Issues

        _ -> NotFound(result.unwrap(path, []))
      }
    }

  #(
    Model(
      authenticated: False,
      route:,
      toasts: dict.new(),
      loading: False,
      wave: Wave(
        id: 1,
        state: WaveHistoric,
        created_at: timestamp.system_time(),
        begin_at: timestamp.system_time(),
        close_at: timestamp.system_time(),
        trainees: [],
      ),
      cases: [],
      user: User(
        id: "system",
        avatar_hash: "",
        display_name: "Unknown",
        role: UnknownUser,
      ),
      issues: [],
      threads: [],
      trainees: [],
      interview_questions: [],
      total_cases: 0,
      total_threads: 0,
      total_issues: 0,
      thread_filter: "",
      threads_open: True,
      threads_closed: True,
      thread_trainee: None,
      view_commands: False,
      searching: False,
      // modal: ThreadIssueModal("a", 0, "c"),
      modal: ClosedModal,
    ),
    effect.batch([
      modem.init(on_url_change),
      get_wave(),
      get_user(),
      route_effects(route),
    ]),
  )
}

type Message {
  OnRouteChange(Route)

  ToastAdded(Toast)

  ToastRemoved(Int)

  // Api returning
  ApiReturnedWave(Result(Wave, rsvp.Error(String)))
  ApiReturnedUser(Result(User, rsvp.Error(String)))
  ApiReturnedThreads(Result(List(ModmailThread), rsvp.Error(String)))
  ApiReturnedThread(Result(ModmailThread, rsvp.Error(String)))
  ApiReturnedCases(Result(List(AthenaCase), rsvp.Error(String)))
  ApiReturnedCase(Result(AthenaCase, rsvp.Error(String)))
  ApiReturnedIssues(Result(List(Issue), rsvp.Error(String)))
  ApiReturnedQuestions(Result(List(String), rsvp.Error(String)))

  ApiReturnedNoAuth

  // User initiated actions
  UserClickedLogin
  UserChangedThreadOpenFilter(Bool)
  UserChangedThreadClosedFilter(Bool)
  UserChangedThreadTraineeFilter(String)
  UserWroteThreadFilter(String)
  UserPromptedThreadIssue(thread_id: String, message_id: Int, mod_id: String)
  UserClosedModal
}

fn on_url_change(uri: uri.Uri) -> Message {
  case uri.path_segments(uri.path) {
    [] | ["threads"] -> OnRouteChange(Threads)

    ["threads", id] -> OnRouteChange(Thread(id, option.None))

    ["stats"] -> OnRouteChange(Stats)

    ["issues"] -> OnRouteChange(Issues)

    ["cases"] -> OnRouteChange(Cases)

    ["questions"] -> OnRouteChange(InterviewQuestions)

    ["api", "auth", "redirect"] -> UserClickedLogin

    ["cases", id] as path ->
      case int.parse(id) {
        Ok(id) -> OnRouteChange(Case(id, option.None))
        _ -> OnRouteChange(NotFound(path))
      }

    path -> OnRouteChange(NotFound(path))
  }
}

fn update(model: Model, message: Message) -> #(Model, Effect(Message)) {
  case message {
    OnRouteChange(route) -> #(Model(..model, route:), route_effects(route))

    // Toasts
    ToastAdded(toast) -> {
      let next_id = dict.size(model.toasts) + 1
      #(
        Model(..model, toasts: dict.insert(model.toasts, next_id, toast)),
        remove_toast_after(next_id, 5000),
      )
    }

    ToastRemoved(id) -> #(
      Model(..model, toasts: dict.drop(model.toasts, [id])),
      effect.none(),
    )

    // API returning messages
    ApiReturnedWave(Ok(wave)) -> #(
      Model(..model, trainees: wave.trainees, wave:),
      effect.none(),
    )

    ApiReturnedWave(Error(err)) -> #(model, rsvp_err_to_toast("wave", err))

    ApiReturnedUser(Ok(user)) -> #(
      Model(..model, user: user, authenticated: True),
      effect.none(),
    )

    ApiReturnedUser(Error(err)) -> #(model, rsvp_err_to_toast("user", err))

    ApiReturnedQuestions(Ok(interview_questions)) -> #(
      Model(..model, interview_questions:),
      effect.none(),
    )

    ApiReturnedQuestions(Error(err)) -> #(
      model,
      rsvp_err_to_toast("questions", err),
    )

    ApiReturnedThreads(Ok(threads)) -> #(
      Model(..model, threads:),
      effect.none(),
    )

    ApiReturnedThreads(Error(err)) -> #(
      model,
      rsvp_err_to_toast("threads", err),
    )

    ApiReturnedThread(Ok(thread)) -> #(
      Model(..model, route: Thread(id: thread.id, content: option.Some(thread))),
      effect.none(),
    )

    ApiReturnedThread(Error(err)) -> #(model, rsvp_err_to_toast("thread", err))

    ApiReturnedCases(Ok(cases)) -> #(Model(..model, cases:), effect.none())

    ApiReturnedCases(Error(err)) -> #(model, rsvp_err_to_toast("cases", err))

    ApiReturnedCase(Ok(athena_case)) -> #(
      Model(
        ..model,
        route: Case(id: athena_case.id, content: option.Some(athena_case)),
      ),
      effect.none(),
    )

    ApiReturnedCase(Error(err)) -> #(model, rsvp_err_to_toast("case", err))

    ApiReturnedIssues(Ok(issues)) -> #(Model(..model, issues:), effect.none())

    ApiReturnedIssues(Error(err)) -> #(model, rsvp_err_to_toast("issues", err))

    ApiReturnedNoAuth -> #(Model(..model, authenticated: False), effect.none())

    // User filtering
    UserChangedThreadOpenFilter(state) -> #(
      Model(..model, threads_open: state),
      effect.none(),
    )

    UserChangedThreadClosedFilter(state) -> #(
      Model(..model, threads_closed: state),
      effect.none(),
    )

    UserChangedThreadTraineeFilter("") -> #(
      Model(..model, thread_trainee: None),
      effect.none(),
    )

    UserChangedThreadTraineeFilter(trainee) -> #(
      Model(..model, thread_trainee: Some(trainee)),
      effect.none(),
    )

    UserWroteThreadFilter(filter) -> #(
      Model(..model, thread_filter: filter),
      effect.none(),
    )

    UserPromptedThreadIssue(thread_id:, message_id:, mod_id:) -> #(
      Model(..model, modal: ThreadIssueModal(thread_id, message_id, mod_id)),
      effect.none(),
    )

    UserClickedLogin -> #(model, push_to_login())

    UserClosedModal -> #(Model(..model, modal: ClosedModal), effect.none())
  }
}

fn view(model: Model) -> Element(Message) {
  case model.authenticated {
    True ->
      html.div(
        [class("grid lg:grid-cols-page h-[100dvh] relative overflow-y-hidden")],
        [
          html.nav(
            [
              class(
                "page-sidebar flex flex-col border-r border-gray-800 h-[100dvh]",
              ),
            ],
            [
              html.details([class("relative")], [
                html.summary([class("sidebar-logo")], [
                  icons.mortarboard([]),
                  html.div([], [
                    html.h3([], [html.text("Academy")]),
                    html.p([class("font-bold text-xs text-gray-400")], [
                      {
                        let #(calendar.Date(year:, month:, ..), _) =
                          timestamp.to_calendar(
                            model.wave.begin_at,
                            calendar.local_offset(),
                          )

                        html.text(
                          calendar.month_to_string(month)
                          <> " "
                          <> int.to_string(year),
                        )
                      },
                    ]),
                  ]),
                  icons.chevron_down([class("size-4 ml-auto")]),
                ]),
                html.ul(
                  [
                    class(
                      "absolute top-full left-3 right-3 bg-gray-900 border border-gray-800 rounded-lg p-1 z-50",
                    ),
                  ],
                  [
                    html.li([], [html.text("Coming soon...")]),
                    // html.li([], [html.button([], [html.text("2026 — June")])]),
                  // html.li([], [html.button([], [html.text("2025 — December")])]),
                  // html.li([], [html.button([], [html.text("2023 — Jure")])]),
                  // html.li([], [html.button([], [html.text("2021 — Septober")])]),
                  ],
                ),
              ]),
              html.ul(
                [],
                list.map(
                  case model.wave.state {
                    WaveInterviews -> [InterviewQuestions]
                    WaveHelper
                      if model.user.role == HelperUser
                      || model.user.role == AdminUser
                    -> [Threads, Cases, Stats, Issues]
                    WaveHelper -> [Threads, Cases, Stats]
                    WaveHistoric -> [Threads, Cases, Stats, Issues]
                  },
                  sidebar_link(model, _),
                ),
              ),
              html.div(
                [
                  class(
                    "mt-auto bg-gray-900 border border-gray-800 rounded-md py-3 px-4 m-3 font-bold text-gray-100",
                  ),
                ],
                [
                  html.text(case model.wave.state {
                    WaveInterviews -> "Interviews"
                    WaveHelper -> "Training"
                    WaveHistoric -> "Historic"
                  }),
                ],
              ),
            ],
          ),

          case model.route {
            Threads | Thread(..) -> threads_sidebar(model)
            Cases | Case(..) -> cases_sidebar(model)
            _ -> element.none()
          },

          html.main(
            [
              class(case model.route {
                Threads | Thread(..) | Cases | Case(..) ->
                  "h-[100dvh] flex flex-col"
                _ -> "lg:col-span-2 h-[100dvh] flex flex-col"
              }),
            ],
            [
              html.header(
                [
                  class(
                    "px-6 h-20 flex items-center justify-between flex-wrap border-b border-gray-900",
                  ),
                ],
                [
                  html.h1([class("font-bold text-xl text-white")], [
                    html.text(case model.route {
                      Threads -> "Threads"
                      Thread(
                        content: option.Some(ModmailThread(user_name:, ..)),
                        ..,
                      ) -> "Thread with " <> user_name
                      Thread(..) -> "Loading thread..."
                      Cases -> "Cases"
                      Case(id:, ..) -> "Case #" <> int.to_string(id)
                      Stats -> "Statistics"
                      Issues -> "Issues"
                      InterviewQuestions -> "Interview Questions"
                      NotFound(_path) -> "Page Not Found"
                    }),
                  ]),
                  html.nav([], [
                    html.details(
                      [class("relative"), attribute.attribute("open", "true")],
                      [
                        html.summary(
                          [
                            class(
                              "flex items-center gap-3 font-semibold cursor-pointer rounded-md py-2 px-3 border border-transparent transition-colors hover:border-gray-800 hover:bg-gray-1000",
                            ),
                          ],
                          [
                            html.img([
                              attribute.src(avatar(model.user)),
                              class("size-7 rounded-full"),
                            ]),
                            html.text(model.user.display_name),
                            icons.chevron_down([class("size-4")]),
                          ],
                        ),
                        html.ul(
                          [
                            class(
                              "absolute top-full right-0 bg-gray-850 rounded-md border border-gray-800 p-1 z-50",
                            ),
                          ],
                          [
                            html.li([], [
                              html.a(
                                [
                                  attribute.href("/api/auth/logout"),
                                  class(
                                    "rounded-md py-1 px-2 flex items-center gap-2 font-semibold",
                                  ),
                                ],
                                [
                                  html.text("Logout"),
                                ],
                              ),
                            ]),
                          ],
                        ),
                      ],
                    ),
                  ]),
                ],
              ),

              case model.route {
                Thread(id:, content: option.None) ->
                  html.div([class("p-6")], [
                    html.div(
                      [
                        attribute.role("alert"),
                        class(
                          "bg-info-bg border border-info-fg text-white p-3 rounded-md",
                        ),
                      ],
                      [
                        html.p([], [
                          html.text(
                            "Loading ModMail thread #" <> id <> " content...",
                          ),
                        ]),
                      ],
                    ),
                  ])

                Thread(content: option.Some(thread), ..) ->
                  thread_view(model, thread)

                Threads ->
                  html.div([class("p-10 text-center text-gray-300 text-xl")], [
                    html.text("Please select a thread"),
                  ])

                Issues -> issues_view(model)

                Cases -> html.div([], [])

                Case(id:, content: option.None) ->
                  html.div([class("p-6")], [
                    html.div(
                      [
                        attribute.role("alert"),
                        class(
                          "bg-info-bg border border-info-fg text-white p-3 rounded-md",
                        ),
                      ],
                      [
                        html.p([], [
                          html.text("Loading case..."),
                        ]),
                      ],
                    ),
                  ])

                Case(content: option.Some(mod_case), ..) ->
                  case_view(model, mod_case)

                Stats -> stats_view(model)

                InterviewQuestions ->
                  questions_view(model.user, model.interview_questions)

                NotFound(_) -> html.div([], [])
              },
            ],
          ),
          keyed.ul(
            [class("fixed top-10 right-4 flex items-end flex-col gap-2")],
            list.map(dict.to_list(model.toasts), fn(combined) {
              let #(key, toast) = combined

              let classes =
                "rounded-lg border border-white/10 py-2 text-sm px-4 font-semibold text-white "
                <> case toast {
                  ToastError(_) -> "bg-rose-900"
                  ToastSuccess(_) -> "bg-emerald-900"
                  ToastWarning(_) -> "bg-orange-900"
                }

              #(
                "toast#" <> int.to_string(key),
                html.li([class(classes)], [html.text(toast.msg)]),
              )
            }),
          ),
          html.div(
            [
              on_direct_click(UserClosedModal),
              class(
                "fixed inset-0 bg-gray-1000/80 flex items-center justify-center transition-opacity",
              ),
              attribute.classes([
                #("opacity-0 pointer-events-none", model.modal == ClosedModal),
              ]),
            ],
            [
              case model.modal {
                ClosedModal -> element.none()
                ThreadIssueModal(thread_id:, message_id:, mod_id:) ->
                  thread_issue_modal(thread_id:, message_id:, mod_id:)
              },
            ],
          ),
        ],
      )
    False ->
      html.div(
        [
          class(
            "h-[100dvh] flex flex-col items-center justify-center py-10 text-center",
          ),
        ],
        [
          icons.mortarboard([class("size-16 fill-ow-orange mb-4")]),
          html.h1([class("font-bold text-3xl text-white")], [
            html.text("You aren't authenticated"),
          ]),
          html.p([class("text-gray-300 max-w-xl mb-6")], [
            html.text(
              "It looks like you are not logged in to Academy. If you are an Overwatch Discord staff member, you can login below.",
            ),
          ]),
          html.a(
            [
              attribute.href("/api/auth/redirect"),
              class(
                "rounded-md py-2 px-4 font-semibold text-white flex items-center gap-2 bg-blurple-500 transition-opacity hover:opacity-80",
              ),
            ],
            [
              icons.discord([class("size-6")]),
              html.text("Login with Discord"),
            ],
          ),
        ],
      )
  }
}

fn sidebar_link(model: Model, route: Route) {
  let #(icon, href, text) = case route {
    Threads -> #(icons.envelope([]), "/threads", "Threads")

    Thread(id:, ..) -> #(icons.inbox([]), "/thread/" <> id, "Thread")

    Cases -> #(icons.cases([]), "/cases", "Cases")

    Case(id:, ..) -> #(icons.inbox([]), "/case" <> int.to_string(id), "Case")

    Stats -> #(icons.graph([]), "/stats", "Statistics")

    Issues -> #(icons.issues([]), "/issues", "Issues")

    InterviewQuestions -> #(
      icons.user_wave([]),
      "/questions",
      "Interview Questions",
    )

    NotFound(_) -> #(icons.inbox([]), "/", "Unknown Link")
  }

  let active = case model.route, route {
    Threads, Threads | Thread(..), Threads -> True
    Cases, Cases | Case(..), Cases -> True
    r1, r2 if r1 == r2 -> True
    _, _ -> False
  }

  html.li([attribute.classes([#("active", active)])], [
    html.a([attribute.href(href)], [
      icon,
      html.text(text),
    ]),
  ])
}

fn threads_sidebar(model: Model) {
  html.aside([class("bg-gray-900 h-[100dvh] flex flex-col")], [
    html.header(
      [
        class("px-4 flex items-center h-20 border-b border-gray-800"),
      ],
      [
        html.form([attribute.action("#"), class("relative w-full")], [
          html.input([
            event.on_input(UserWroteThreadFilter),
            attribute.placeholder("Search thread user..."),
            class(
              "border border-gray-800 bg-gray-950 rounded-md py-1.5 px-4 w-full transition-colors outline-none",
            ),
          ]),
        ]),
      ],
    ),
    html.nav([class("filters p-4 flex gap-2 flex-wrap")], [
      html.button(
        [
          event.on_click(UserChangedThreadOpenFilter(!model.threads_open)),
          class(
            "rounded-sm py-1 px-3 font-semibold flex items-center gap-2 transition-all cursor-pointer hover:opacity-80 "
            <> case model.threads_open {
              True -> "bg-gray-800 text-gray-200"
              False -> "bg-gray-950 text-gray-400"
            },
          ),
        ],
        [
          icons.checkmark([
            class("transition-all size-5"),
            attribute.classes([#("-mr-6 opacity-0", !model.threads_open)]),
          ]),
          html.text("Open"),
        ],
      ),
      html.button(
        [
          event.on_click(UserChangedThreadClosedFilter(!model.threads_closed)),
          class(
            "bg-gray-800 rounded-sm py-1 px-3 font-semibold flex items-center gap-1 transition-all cursor-pointer hover:opacity-80 "
            <> case model.threads_closed {
              True -> "bg-gray-800 text-gray-200"
              False -> "bg-gray-950 text-gray-400"
            },
          ),
        ],
        [
          icons.checkmark([
            class("transition-all size-5"),
            attribute.classes([#("-mr-6 opacity-0", !model.threads_closed)]),
          ]),
          html.text("Closed"),
        ],
      ),
      html.select(
        [
          event.on_change(UserChangedThreadTraineeFilter),
          class(
            "bg-gray-800 rounded-sm py-1 px-3 font-semibold flex items-center gap-2 flex-1 transition-opacity cursor-pointer hover:opacity-80",
          ),
        ],
        [
          html.option([attribute.value("")], "Any Trainee"),
          ..list.map(model.trainees, fn(trainee) {
            html.option([attribute.value(trainee.id)], trainee.username)
          })
        ],
      ),
    ]),

    keyed.ul(
      [class("grid gap-2 px-4 overflow-y-auto flex-1 pb-6")],
      case filtered_threads(model) {
        [] -> [
          #(
            "0",
            html.li([class("text-center leading-loose p-6")], [
              html.text("Sorry, no threads match your criteria."),
            ]),
          ),
        ]
        threads ->
          list.map(threads, fn(thread) {
            #(
              thread.id,
              html.li([class("group")], [
                html.a(
                  [
                    attribute.href("/threads/" <> thread.id),
                    class("block p-5 rounded-md border border-l-4"),
                    attribute.classes([
                      #(
                        "bg-gray-800 border-gray-750 text-gray-200",
                        thread.status == ThreadOpen,
                      ),
                      #(
                        "bg-gray-950 border-gray-800 text-gray-300",
                        thread.status == ThreadClosed,
                      ),
                    ]),
                  ],
                  [
                    html.h4(
                      [class("font-semibold mb-2 flex items-center gap-1.5")],
                      [
                        icons.hashtag([class("size-5 text-gray-400")]),
                        html.text(thread.user_name),
                      ],
                    ),
                    html.dl(
                      [
                        class(
                          "text-gray-200 text-sm flex items-end flex-wrap gap-4 font-semibold",
                        ),
                      ],
                      [
                        html.dd(
                          [
                            class("flex items-center gap-1.5"),
                            attribute.data("tooltip", "User Messages"),
                          ],
                          [
                            icons.arrow_in([class("size-5 text-gray-400")]),
                            html.text(int.to_string(thread.inbound_messages)),
                          ],
                        ),
                        html.dd(
                          [
                            class("flex items-center gap-1.5"),
                            attribute.data("tooltip", "Replies"),
                          ],
                          [
                            icons.arrow_out([
                              class("size-5 text-gray-400"),
                            ]),
                            html.text(int.to_string(thread.outbound_messages)),
                          ],
                        ),
                        html.dd(
                          [
                            class("flex items-center gap-1.5 opacity-90"),
                            attribute.data("tooltip", "Internal Chat"),
                          ],
                          [
                            icons.message_bubbles([
                              class("size-5 text-gray-400"),
                            ]),
                            html.text(int.to_string(thread.chat_messages)),
                          ],
                        ),
                        html.dt([class("ml-auto")], [
                          html.div(
                            [class("flex")],
                            list.map(thread.participants, fn(snowflake) {
                              html.img([
                                attribute.src(""),
                                class(
                                  "size-7 rounded-full bg-blue-400 border border-gray-800 not-last:-mr-2",
                                ),
                              ])
                            }),
                          ),
                        ]),
                      ],
                    ),
                  ],
                ),
              ]),
            )
          })
      },
    ),
  ])
}

fn cases_sidebar(model: Model) {
  html.aside([class("bg-gray-900 h-[100dvh] flex flex-col")], [
    html.header(
      [
        class("px-4 flex items-center h-20 border-b border-gray-800"),
      ],
      [
        html.form([attribute.action("#"), class("relative w-full")], [
          html.input([
            event.on_input(UserWroteThreadFilter),
            attribute.placeholder("Search case user..."),
            class(
              "border border-gray-800 bg-gray-950 rounded-md py-1.5 px-4 w-full transition-colors outline-none",
            ),
          ]),
        ]),
      ],
    ),
    html.nav([class("filters p-4 flex gap-2 flex-wrap")], [
      html.select(
        [
          event.on_change(UserChangedThreadTraineeFilter),
          class(
            "bg-gray-800 rounded-sm py-1 px-3 font-semibold flex items-center gap-2 flex-1 transition-opacity cursor-pointer hover:opacity-80",
          ),
        ],
        [
          html.option([attribute.value("")], "Any Case Type"),
          html.option([attribute.value("warns")], "Warns"),
          html.option([attribute.value("mutes")], "Mutes"),
          html.option([attribute.value("bans")], "Bans"),
        ],
      ),

      html.select(
        [
          event.on_change(UserChangedThreadTraineeFilter),
          class(
            "bg-gray-800 rounded-sm py-1 px-3 font-semibold flex items-center gap-2 flex-1 transition-opacity cursor-pointer hover:opacity-80",
          ),
        ],
        [
          html.option([attribute.value("")], "Any Trainee"),
          ..list.map(model.trainees, fn(trainee) {
            html.option([attribute.value(trainee.id)], trainee.username)
          })
        ],
      ),
    ]),

    keyed.ul(
      [class("grid items-start gap-2 px-4 overflow-y-auto flex-1 pb-6")],
      list.map(model.cases, fn(mod_case) {
        #(
          "case#" <> int.to_string(mod_case.id),
          html.li([], [
            html.a(
              [
                attribute.href("/cases/" <> int.to_string(mod_case.id)),
                class(
                  "border border-gray-750 bg-gray-800 rounded border-l-4 py-2 px-4 grid gap-2 "
                  <> case mod_case.kind {
                    CaseBan -> "border-l-case-ban "
                    CaseUnban -> "border-l-case-unban "
                    CaseNote -> "border-l-case-note "
                    CaseWarn -> "border-l-case-warn "
                    CaseKick -> "border-l-case-kick "
                    CaseMute -> "border-l-case-mute "
                    CaseUnmute -> "border-l-case-unmute "
                    CaseDeleted -> "border-l-case-deleted "
                    CaseSoftban -> "border-l-case-softban "
                  }
                  <> case model.route {
                    Case(id:, ..) if id == mod_case.id ->
                      "bg-gray-800 border-gray-750 text-gray-200"
                    _ -> "bg-gray-950 border-gray-800 text-gray-300"
                  },
                ),
              ],
              [
                html.div(
                  [class("flex items-center justify-between gap-3 flex-wrap")],
                  [
                    html.h3([class("font-semibold text-white text-base")], [
                      html.text(case mod_case.kind {
                        CaseBan -> "Ban"
                        CaseUnban -> "Unban"
                        CaseNote -> "Note"
                        CaseWarn -> "Warn"
                        CaseKick -> "Kick"
                        CaseMute -> "Mute"
                        CaseUnmute -> "Unmute"
                        CaseDeleted -> "Deleted"
                        CaseSoftban -> "Softban"
                      }),
                    ]),
                    html.span(
                      [
                        class(
                          "px-2 pt-0.5 pb-1 rounded-md bg-tag-bg text-tag-fg leading-none",
                        ),
                      ],
                      [html.text(mod_case.mod_name)],
                    ),
                  ],
                ),

                html.p([], [
                  html.text(mod_case.actioned_user_name),
                ]),
                html.footer([class("text-sm")], [
                  html.p([], [
                    html.text("Created "),
                    html.time([class("bg-gray-750 rounded-sm px-1")], [
                      html.text("June 24th, 2026"),
                    ]),
                  ]),
                ]),
              ],
            ),
          ]),
        )
      }),
    ),
  ])
}

fn filtered_threads(model: Model) -> List(ModmailThread) {
  model.threads
  |> list.filter(fn(thread) {
    case model.thread_filter {
      "" -> True
      _ ->
        string.contains(
          string.lowercase(thread.user_name),
          string.lowercase(model.thread_filter),
        )
    }
  })
}

fn questions_view(user: User, questions: List(String)) {
  html.div([class("p-6 bg-gray-900")], [
    html.p(
      [
        class(
          "text-white border border-info-fg bg-info-bg rounded-md px-4 py-2 flex items-center gap-2 mb-8",
        ),
      ],
      [
        icons.info_circle([class("size-5 text-info-fg")]),
        html.text(
          "Clicking the question text will automatically copy-paste it into your clipboard, prefixed with !ar.",
        ),
        case user.role {
          AdminUser ->
            html.a(
              [
                attribute.href("/edit-questions"),
                class(
                  "ml-auto bg-gray-800 rounded-md py-1 px-4 transition-colors cursor-pointer hover:bg-gray-750",
                ),
              ],
              [
                html.text("Edit"),
              ],
            )
          _ -> element.none()
        },
      ],
    ),
    case questions {
      [] ->
        html.div(
          [
            attribute.role("alert"),
            class(
              "text-white border border-orange-300 bg-orange-500/10 rounded-md px-4 py-2 flex items-center gap-2 mb-8",
            ),
          ],
          [
            icons.info_circle([class("size-5 text-orange-300")]),
            html.text("Loading interview questions..."),
          ],
        )
      _ ->
        keyed.ul(
          [class("grid gap-4")],
          list.index_map(questions, fn(question, i) {
            #(
              "question-" <> int.to_string(i),
              html.li(
                [
                  class(
                    "flex items-center px-4 py-5 bg-gray-800 border border-gray-750 rounded-lg",
                  ),
                ],
                [
                  html.input([
                    attribute.type_("checkbox"),
                    class("opacity-0 absolute"),
                  ]),
                  html.div(
                    [
                      attribute.class(
                        "bg-gray-800 border border-gray-750 p-1 rounded-md mr-4",
                      ),
                    ],
                    [
                      svg.svg(
                        [
                          attribute.class("size-5"),
                          attribute("stroke-linejoin", "round"),
                          attribute("stroke-linecap", "round"),
                          attribute("stroke-width", "3"),
                          attribute("stroke", "currentColor"),
                          attribute("fill", "none"),
                          attribute("viewBox", "0 0 24 24"),
                          attribute("xmlns", "http://www.w3.org/2000/svg"),
                        ],
                        [svg.path([attribute("d", "M20 6 9 17l-5-5")])],
                      ),
                    ],
                  ),
                  // html.p([attribute.class("question-num")], [
                  //   html.text("#" <> int.to_string(i + 1)),
                  // ]),
                  html.div([attribute.class("questions")], [
                    html.p(
                      [
                        attribute.class(
                          "question hover:underline cursor-pointer",
                        ),
                        attribute.attribute(
                          "onclick",
                          "window.navigator.clipboard.writeText(this.textContent)",
                        ),
                      ],
                      [
                        html.text(question),
                      ],
                    ),
                  ]),
                ],
              ),
            )
          }),
        )
    },
  ])
}

fn stats_view(model: Model) {
  html.main([class("p-6 bg-gray-900 h-full")], [
    html.section([class("grid lg:grid-cols-3 gap-8")], [
      html.article(
        [
          class(
            "bg-gray-800 border border-gray-750 rounded-xl p-8 flex items-center gap-8",
          ),
        ],
        [
          html.figure([class("p-4 rounded-lg bg-orange-500/10")], [
            icons.envelope([class("size-8 text-orange-300")]),
          ]),
          html.div([], [
            html.h1([class("text-3xl font-extrabold text-white")], [
              html.text(int.to_string(model.total_threads)),
            ]),
            html.h3([class("text-gray-300")], [
              html.text("Threads"),
            ]),
          ]),
        ],
      ),
      html.article(
        [
          class(
            "bg-gray-800 border border-gray-750 rounded-xl p-8 flex items-center gap-8",
          ),
        ],
        [
          html.figure([class("p-4 rounded-lg bg-blue-500/10")], [
            icons.cases([class("size-8 text-blue-300")]),
          ]),
          html.div([], [
            html.h1([class("text-3xl font-extrabold text-white")], [
              html.text(int.to_string(model.total_cases)),
            ]),
            html.h3([class("text-gray-300")], [
              html.text("Cases"),
            ]),
          ]),
        ],
      ),
      html.article(
        [
          class(
            "bg-gray-800 border border-gray-750 rounded-xl p-8 flex items-center gap-8",
          ),
        ],
        [
          html.figure([class("p-4 rounded-lg bg-rose-500/10")], [
            icons.issues([class("size-8 text-rose-300")]),
          ]),
          html.div([], [
            html.h1([class("text-3xl font-extrabold text-white")], [
              html.text(int.to_string(model.total_issues)),
            ]),
            html.h3([class("text-gray-300")], [
              html.text("Issues"),
            ]),
          ]),
        ],
      ),

      html.header([class("")], [
        html.h3([class("text-xl font-bold text-white")], [
          html.text("Trainees"),
        ]),
      ]),

      keyed.ul(
        [class("grid lg:col-span-3 gap-4")],
        list.map(model.trainees, fn(trainee) {
          #(
            "trainee#" <> trainee.id,
            html.li(
              [
                class(
                  "rounded flex items-center gap-4 p-2 hover:bg-gray-800 transition-colors",
                ),
              ],
              [
                html.figure(
                  [
                    class("size-14 rounded-full bg-black bg-cover bg-center"),
                    attribute.style("background-image", "url(" <> "" <> ")"),
                  ],
                  [],
                ),
                html.div([], [
                  html.h3([class("font-semibold text-lg text-ow-mod")], [
                    html.text(trainee.username),
                  ]),
                  html.p([class("text-gray-300")], [
                    html.text(
                      int.to_string(trainee.message_count)
                      <> " messages, "
                      <> int.to_string(trainee.thread_participation_count)
                      <> " threads participated in, "
                      <> int.to_string(trainee.case_count)
                      <> " cases",
                    ),
                  ]),
                ]),
              ],
            ),
          )
        }),
      ),
    ]),
  ])
}

fn thread_view(model: Model, thread: ModmailThread) {
  html.div([class("py-6 block h-full overflow-y-auto")], [
    html.header([class("mx-6 pb-6 mb-6 border-b border-gray-800")], [
      html.h2([class("font-semibold text-white text-3xl mb-4")], [
        html.text(
          case thread.status {
            ThreadOpen -> "Open"
            ThreadClosed -> "Closed"
            ThreadSuspended -> "Suspended"
          }
          <> " thread with "
          <> thread.user_name,
        ),
      ]),
      html.p([class("mb-4")], [
        html.code(
          [
            class(
              "font-monospace border border-gray-800 rounded-md inline-block px-3",
            ),
          ],
          [html.text("user #" <> thread.user_id)],
        ),
      ]),
      html.div([class("flex justify-between gap-6 items-center")], [
        html.ul(
          [class("flex flex-wrap gap-2")],
          list.map(thread.roles, fn(role) {
            html.li(
              [
                class(
                  "rounded-full py-1 px-3 font-semibold border text-sm "
                  <> case role {
                    "Muted" -> "bg-closed/10 border-closed text-red-100"
                    "Admin" | "Administrator" ->
                      "bg-ow-admin/10 border-ow-admin text-blue-100"
                    "Mod" | "Moderator" ->
                      "bg-ow-mod/10 border-ow-mod text-orange-100"
                    _ -> "bg-gray-800 border-gray-750 text-gray-100"
                  },
                ),
              ],
              [html.text(role)],
            )
          }),
        ),
        html.ul(
          [],
          list.map(thread.participants, fn(participant) {
            html.li([], [
              html.img([
                class("rounded-full size-8"),
                attribute.src(""),
                attribute.alt("Participant avatar"),
              ]),
            ])
          }),
        ),
      ]),
    ]),
    keyed.ul(
      [class("grid")],
      list.index_map(thread.messages, fn(message, i) {
        #(
          "msg#" <> message.user_id <> int.to_string(i),
          html.li(
            [
              class(
                "flex gap-3 px-8 py-3 transition-colors hover:bg-gray-900 relative group "
                <> case message.kind {
                  InternalMsg -> "bg-gray-900/50"
                  CommandMsg if !model.view_commands -> "hidden"
                  _ -> ""
                },
              ),
            ],
            [
              html.figure([], [
                html.img([
                  class("size-11 rounded-full bg-black"),
                  attribute.alt(message.user_name <> "'s Avatar"),
                  attribute.src(case message.kind {
                    IncomingMsg -> ""
                    _ -> ""
                    // avatar(message.user_id)
                  }),
                ]),
              ]),
              html.section([class("flex-1")], [
                html.h4(
                  [
                    class(
                      "flex items-center gap-2 font-semibold "
                      <> case message.role {
                        ModRole | TraineeRole -> "text-ow-mod"
                        HelperRole -> "text-ow-helper"
                        AdminRole -> "text-ow-admin"
                        SystemRole -> ""
                        UserRole -> "text-white"
                      },
                    ),
                  ],
                  [
                    html.text(case message.kind {
                      SystemMsg -> "ModMail"
                      _ -> message.user_name
                    }),
                    html.span(
                      [
                        class(
                          "text-xs text-gray-200 bg-gray-800 leading-none pt-0.5 pb-1 rounded px-1.5 uppercase",
                        ),
                      ],
                      [
                        html.text(case message.kind {
                          InternalMsg -> "Internal"
                          IncomingMsg -> "From User"
                          OutgoingMsg -> "To User"
                          SystemMsg -> "System"
                          CommandMsg -> "Command"
                        }),
                      ],
                    ),
                  ],
                ),
                element.unsafe_raw_html(
                  "",
                  "article",
                  [class("message-content")],
                  message.body,
                ),
                case message.attachments {
                  [] -> element.none()
                  _ ->
                    html.footer(
                      [class("mt-1")],
                      list.map(message.attachments, fn(attachment) {
                        case string.reverse(attachment) {
                          "gnp." <> _ ->
                            html.img([
                              attribute.src(attachment),
                              attribute.alt("Modmail Embedded Image"),
                              class("max-h-96 max-w-96 rounded-md"),
                            ])
                          _ -> element.none()
                        }
                      }),
                    )
                },
              ]),
              html.button(
                [
                  event.on_click(UserPromptedThreadIssue(
                    thread.id,
                    message.id,
                    message.user_id,
                  )),
                  attribute.data("tooltip", "Raise Issue"),
                  class(
                    "absolute top-0 right-8 shadow-sm bg-gray-900 border border-gray-800 rounded-sm p-2 -translate-y-1/2 opacity-0 pointer-events-none transition-opacity group-hover:opacity-100 group-hover:pointer-events-auto cursor-pointer",
                  ),
                ],
                [icons.issues([class("size-4")])],
              ),
            ],
          ),
        )
      }),
    ),
  ])
}

fn case_view(model: Model, mod_case: AthenaCase) {
  html.div([class("py-6 px-6 block h-full overflow-y-auto")], [
    html.header([], []),
    keyed.ul(
      [class("grid gap-4")],
      list.index_map(mod_case.notes, fn(note, i) {
        #(
          "note#" <> int.to_string(i),
          html.li([], [element.unsafe_raw_html("", "article", [], note.body)]),
        )
      }),
    ),
  ])
}

fn issues_view(model: Model) {
  html.div([class("p-6 bg-gray-900")], [
    html.header([class("flex items-center justify-between gap-4 mb-8")], [
      html.h2([class("text-2xl text-gray-100")], [
        html.text(int.to_string(list.length(model.issues)) <> " issues found"),
      ]),
    ]),
    html.section([], [
      keyed.ul(
        [class("grid gap-4")],
        list.map(model.issues, fn(issue) {
          #(
            "issue#" <> int.to_string(issue.id),
            html.li([], [
              html.a(
                [
                  attribute.href("/issues/" <> int.to_string(issue.id)),
                  class(
                    "block border border-gray-750 bg-gray-800 rounded border-l-3 py-4 px-5",
                  ),
                  attribute.classes([
                    #("border-l-case-blue", issue.status == "review"),
                    #("border-l-gray-600", issue.status == "handled"),
                    #("border-l-red-500", issue.status == "focus"),
                  ]),
                ],
                [
                  html.h3([class("font-bold text-white")], [
                    html.text(issue.category),
                  ]),
                  html.p([], [
                    html.text(issue.reason),
                  ]),
                  // html.ul([class("flex gap-5 flex-wrap")], [
                //   html.li([], [
                //     html.h5([class("text-white font-semibold")], [
                //       html.text("Reported by"),
                //     ]),
                //     html.p([], [
                //       html.button(
                //         [
                //           class(
                //             "px-1 py-0.5 rounded-md bg-tag-bg text-tag-fg leading-none",
                //           ),
                //         ],
                //         [html.text("@graphiteisaac")],
                //       ),
                //     ]),
                //   ]),
                //
                //   html.li([], [
                //     html.h5([class("text-white font-semibold")], [
                //       html.text("Trainee"),
                //     ]),
                //     html.p([], [
                //       html.button(
                //         [
                //           class(
                //             "px-1 py-0.5 rounded-md bg-tag-bg text-tag-fg leading-none",
                //           ),
                //         ],
                //         [html.text("@poopsocket")],
                //       ),
                //     ]),
                //   ]),
                // ]),
                // html.footer([class("text-sm")], [
                //   html.p([], [
                //     html.text("Created "),
                //     html.time([class("bg-gray-750 rounded-sm px-1")], [
                //       html.text("June 24th, 2026"),
                //     ]),
                //   ]),
                // ]),
                ],
              ),
            ]),
          )
        }),
      ),
    ]),
  ])
}

fn thread_issue_modal(
  thread_id thread_id: String,
  message_id message_id: Int,
  mod_id mod_id: String,
) {
  html.div(
    [
      class("bg-gray-900 rounded-xl max-w-xl w-full shadow-xl"),
    ],
    [
      html.header(
        [
          class("p-5 lg:p-8 flex items-center gap-3 flex-wrap justify-between"),
        ],
        [
          html.h4([class("font-semibold text-white text-xl")], [
            html.text("Raise an issue"),
          ]),
          html.button(
            [
              event.on_click(UserClosedModal),
              class(
                "cursor-pointer p-2 cursor-pointer rounded-md hover:bg-gray-800",
              ),
            ],
            [icons.x([class("size-4")])],
          ),
        ],
      ),
      html.form(
        [
          class("p-5 lg:p-8 pt-0 lg:pt-0"),
          attribute.method("post"),
        ],
        [
          html.input([
            attribute.type_("hidden"),
            attribute.name("thread_id"),
            attribute.value(thread_id),
          ]),
          html.input([
            attribute.type_("hidden"),
            attribute.name("message_id"),
            attribute.value(int.to_string(message_id)),
          ]),
          html.input([
            attribute.type_("hidden"),
            attribute.name("mod_id"),
            attribute.value(mod_id),
          ]),
          html.div([class("form-row")], [
            html.label([attribute.for("concern")], [
              html.text("Categorize your concern"),
            ]),
            html.select(
              [
                attribute.id("concern"),
                attribute.name("concern"),
              ],
              [
                html.option(
                  [attribute.value("bad_response")],
                  "Poorly communicated response",
                ),
                html.option(
                  [attribute.value("against_policy")],
                  "Against our policies",
                ),
                html.option(
                  [attribute.value("oversharing")],
                  "Oversharing information",
                ),
                html.option([attribute.value("argumentative")], "Argumentative"),
              ],
            ),
          ]),

          html.div([class("form-row")], [
            html.label([attribute.for("thoughts")], [
              html.text("Briefly describe your thoughts"),
            ]),
            html.textarea(
              [
                attribute.id("thoughts"),
                attribute.name("thoughts"),
                attribute.rows(3),
              ],
              "",
            ),
          ]),

          html.div([class("form-row submission-row")], [
            html.button([attribute.type_("submit")], [
              html.text("Finish Raising"),
            ]),
          ]),
        ],
      ),
    ],
  )
}

//
// Data functions
//

fn get_wave() -> Effect(Message) {
  let handler = rsvp.expect_json(wave_decoder(), ApiReturnedWave)
  rsvp.get(base_url <> "/api/wave", handler)
}

fn get_user() -> Effect(Message) {
  let handler = rsvp.expect_json(user_decoder(), ApiReturnedUser)
  rsvp.get(base_url <> "/api/auth/me", handler)
}

fn get_questions() -> Effect(Message) {
  let handler =
    rsvp.expect_json(decode.list(decode.string), ApiReturnedQuestions)
  rsvp.get(base_url <> "/api/questions", handler)
}

fn get_threads() -> Effect(Message) {
  let handler =
    rsvp.expect_json(
      {
        use threads <- decode.field(
          "threads",
          decode.list(modmail_thread_decoder()),
        )
        decode.success(threads)
      },
      ApiReturnedThreads,
    )

  rsvp.get(base_url <> "/api/threads", handler)
}

fn get_thread(id: String) -> Effect(Message) {
  let handler = rsvp.expect_json(modmail_thread_decoder(), ApiReturnedThread)
  rsvp.get(base_url <> "/api/threads/" <> id, handler)
}

fn get_cases() -> Effect(Message) {
  let handler =
    rsvp.expect_json(
      {
        use threads <- decode.field("cases", decode.list(athena_case_decoder()))
        decode.success(threads)
      },
      ApiReturnedCases,
    )

  rsvp.get(base_url <> "/api/cases", handler)
}

fn get_case(id: Int) -> Effect(Message) {
  let handler = rsvp.expect_json(athena_case_decoder(), ApiReturnedCase)
  rsvp.get(base_url <> "/api/cases/" <> int.to_string(id), handler)
}

fn get_issues() -> Effect(Message) {
  let handler =
    rsvp.expect_json(decode.list(issue_decoder()), ApiReturnedIssues)
  rsvp.get(base_url <> "/api/issues", handler)
}

//
// Custom effects
//

fn set_title(route: Route) -> Effect(Message) {
  let page_title = case route {
    Threads -> "Threads"
    Thread(id:, ..) -> "Thread #" <> id
    Cases -> "Cases"
    Case(id:, ..) -> "Case #" <> int.to_string(id)
    Stats -> "Statistics"
    Issues -> "Issues"
    InterviewQuestions -> "Interview Questions"
    NotFound(..) -> "Page not found"
  }

  use _ <- effect.from
  browser.set_page_title(page_title <> " ・ Academy")
}

fn remove_toast_after(id: Int, delay: Int) -> Effect(Message) {
  use dispatch <- effect.from
  browser.set_timeout(delay, fn() { dispatch(ToastRemoved(id)) })
}

fn push_to_login() {
  use _ <- effect.from
  browser.push_to_url("/api/auth/redirect")
}

// Custom events

fn on_direct_click(msg: msg) -> attribute.Attribute(msg) {
  let decoder = {
    use target <- decode.field("target", decode.dynamic)
    use current <- decode.field("currentTarget", decode.dynamic)

    case browser.is_same_node(target, current) {
      True -> decode.success(msg)
      False -> decode.failure(msg, "targets did not match")
    }
  }

  event.on("click", decoder)
}

// Utils

fn route_effects(route: Route) -> Effect(Message) {
  let data_effects = case route {
    Threads -> [get_threads()]
    Thread(id:, content: option.None) -> [get_threads(), get_thread(id)]
    Cases -> [get_cases()]
    Case(id:, content: option.None) -> [get_cases(), get_case(id)]
    Stats -> []
    Issues -> [get_issues()]
    InterviewQuestions -> [
      get_questions(),
    ]

    _ -> []
  }

  effect.batch([set_title(route), ..data_effects])
}

fn rsvp_err_to_toast(
  resource: String,
  err: rsvp.Error(String),
) -> Effect(Message) {
  echo #(resource, err)

  use dispatch <- effect.from

  case err {
    rsvp.BadBody ->
      ToastAdded(ToastError("The response body could not be decoded"))
    rsvp.BadUrl(_) ->
      ToastAdded(ToastError("The provided URL was badly formed"))
    rsvp.HttpError(response.Response(status: 401, ..)) -> ApiReturnedNoAuth
    rsvp.HttpError(response.Response(status:, ..)) ->
      ToastAdded(ToastError(
        "Couldn't get " <> resource <> ", HTTP status " <> int.to_string(status),
      ))
    rsvp.JsonError(_) ->
      ToastAdded(ToastError("Could not decode JSON from response"))
    rsvp.NetworkError -> ToastAdded(ToastError("A network error has occurred"))
    rsvp.UnhandledResponse(_) ->
      ToastAdded(ToastError("A response wasn't handled properly"))
  }
  |> dispatch
}

fn avatar(user: User) {
  base_url <> "/api/avatar/" <> user.id <> "/" <> user.avatar_hash <> ".png"
}
