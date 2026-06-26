import lustre/attribute.{attribute}
import lustre/element.{type Element}
import lustre/element/svg

pub fn inbox(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "currentColor"),
      attribute.role("img"),
      attribute("viewBox", "0 0 24 24"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute.class(""),
        attribute("clip-rule", "evenodd"),
        attribute(
          "d",
          "M5 2a3 3 0 0 0-3 3v14a3 3 0 0 0 3 3h14a3 3 0 0 0 3-3V5a3 3 0 0 0-3-3H5ZM4 5.5C4 4.67 4.67 4 5.5 4h13c.83 0 1.5.67 1.5 1.5v6c0 .83-.67 1.5-1.5 1.5h-2.65c-.5 0-.85.5-.85 1a3 3 0 1 1-6 0c0-.5-.35-1-.85-1H5.5A1.5 1.5 0 0 1 4 11.5v-6Z",
        ),
        attribute("fill-rule", "evenodd"),
      ]),
    ],
  )
}

pub fn graph(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "currentColor"),
      attribute.role("img"),
      attribute("viewBox", "0 0 24 24"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute("clip-rule", "evenodd"),
        attribute(
          "d",
          "M2 19V5a3 3 0 0 1 3-3h14a3 3 0 0 1 3 3v14a3 3 0 0 1-3 3H5a3 3 0 0 1-3-3Zm16-9.59V13a1 1 0 1 0 2 0V7a1 1 0 0 0-1-1h-6a1 1 0 1 0 0 2h3.59l-5.09 5.09-1.8-1.8a1 1 0 0 0-1.4 0l-4 4a1 1 0 1 0 1.4 1.42L9 13.4l1.8 1.8a1 1 0 0 0 1.4 0L18 9.4Z",
        ),
        attribute("fill-rule", "evenodd"),
      ]),
    ],
  )
}

pub fn issues(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "currentColor"),
      attribute("viewBox", "0 0 24 24"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute(
          "d",
          "M3 1a1 1 0 0 1 1 1v.82l8.67-1.45A2 2 0 0 1 15 3.35v1.47l5.67-.95A2 2 0 0 1 23 5.85v7.3a2 2 0 0 1-1.67 1.98l-9 1.5a2 2 0 0 1-1.78-.6c-.2-.21-.08-.54.18-.68a5.01 5.01 0 0 0 1.94-1.94c.18-.32-.1-.66-.46-.6L4 14.18V21a1 1 0 1 1-2 0V2a1 1 0 0 1 1-1Z",
        ),
      ]),
    ],
  )
}

pub fn users(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "none"),
      attribute("stroke", "currentColor"),
      attribute("stroke-linecap", "round"),
      attribute("stroke-linejoin", "round"),
      attribute("stroke-width", "2"),
      attribute("viewBox", "0 0 24 24"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([attribute("d", "M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2")]),
      svg.path([attribute("d", "M16 3.128a4 4 0 0 1 0 7.744")]),
      svg.path([attribute("d", "M22 21v-2a4 4 0 0 0-3-3.87")]),
      svg.circle([
        attribute("cx", "9"),
        attribute("cy", "7"),
        attribute("r", "4"),
      ]),
    ],
  )
}

pub fn school(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "none"),
      attribute("stroke", "currentColor"),
      attribute("stroke-linecap", "round"),
      attribute("stroke-linejoin", "round"),
      attribute("stroke-width", "2"),
      attribute("viewBox", "0 0 24 24"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([attribute("d", "M14 21v-3a2 2 0 0 0-4 0v3")]),
      svg.path([attribute("d", "M18 4.933V21")]),
      svg.path([attribute("d", "m4 6 7.106-3.79a2 2 0 0 1 1.788 0L20 6")]),
      svg.path([
        attribute(
          "d",
          "m6 11-3.52 2.147a1 1 0 0 0-.48.854V19a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-5a1 1 0 0 0-.48-.853L18 11",
        ),
      ]),
      svg.path([attribute("d", "M6 4.933V21")]),
      svg.circle([
        attribute("cx", "12"),
        attribute("cy", "9"),
        attribute("r", "2"),
      ]),
    ],
  )
}

pub fn shield(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute.aria_label("Shield Icon"),
      attribute("fill", "currentColor"),
      attribute.role("img"),
      attribute("viewBox", "0 0 20 20"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute(
          "d",
          "M15.3296 0H4.6586C4.49194 2.16758 2.65782 3.83491 0.490234 3.83491V4.8353C0.490234 9.75393 2.82462 14.3392 6.90952 17.5905L9.99416 20.008L13.0787 17.5905C17.1637 14.4224 19.498 9.75393 19.498 4.8353V3.83491C17.3304 3.83491 15.5797 2.16758 15.3296 0ZM8.07671 14.4224C5.57566 12.4216 4.07501 9.58726 4.07501 6.50262V5.91902C5.40887 5.91902 6.57606 4.91863 6.65939 3.58477H9.99416V16.0064L8.07671 14.4224Z",
        ),
      ]),
    ],
  )
}

pub fn mortarboard(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "currentColor"),
      attribute("viewBox", "0 0 502 526"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute("clip-rule", "evenodd"),
        attribute(
          "d",
          "M246.807 97.6691C249.247 96.7766 251.924 96.7766 254.364 97.6691L494.74 185.601C504.368 189.123 504.368 202.74 494.74 206.262L255.155 293.904C252.715 294.797 250.038 294.797 247.598 293.904L105.643 241.976C96.5447 255.321 91.0945 270.702 87.9492 286.144C85.812 296.637 84.792 306.899 84.3711 316.178C93.0748 321.605 98.9813 332.158 98.9814 344.287C98.9813 355.051 94.3297 364.573 87.2002 370.373L99.0547 483.453C99.5493 488.175 95.8464 492.287 91.0986 492.287H50.8633C46.1154 492.287 42.4126 488.175 42.9072 483.453L54.7617 370.374C47.632 364.574 42.9816 355.051 42.9814 344.287C42.9816 332.761 48.3143 322.659 56.3135 317.025C56.7251 306.052 57.8715 293.523 60.5127 280.556C63.7355 264.734 69.2721 247.816 78.7109 232.124L7.22071 205.973C-2.40684 202.451 -2.40696 188.834 7.22071 185.312L246.807 97.6691ZM236.604 323.956C246.143 327.445 256.609 327.446 266.148 323.956L385.527 280.286H385.98L389.84 359.412C389.931 360.365 389.98 361.323 389.98 362.286C389.98 397.632 327.748 426.286 250.98 426.286C174.213 426.286 111.98 397.632 111.98 362.286C111.98 361.323 112.029 360.365 112.12 359.412L115.98 280.286H117.225L236.604 323.956Z",
        ),
        attribute("fill-rule", "evenodd"),
      ]),
    ],
  )
}

pub fn envelope(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "currentColor"),
      attribute.role("img"),
      attribute("viewBox", "0 0 24 24"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute(
          "d",
          "M1.16 5.02c-.1.28.04.58.29.74l10.27 6.85a.5.5 0 0 0 .56 0l10.27-6.85c.25-.16.38-.46.29-.74A3 3 0 0 0 20 3H4a3 3 0 0 0-2.84 2.02Z",
        ),
      ]),
      svg.path([
        attribute(
          "d",
          "M23 8.8a.5.5 0 0 0-.78-.41l-9.53 6.35c-.42.28-.96.28-1.38 0L1.78 8.39A.5.5 0 0 0 1 8.8V18a3 3 0 0 0 3 3h16a3 3 0 0 0 3-3V8.8Z",
        ),
      ]),
    ],
  )
}

pub fn cases(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "currentColor"),
      attribute.role("img"),
      attribute("viewBox", "0 0 24 24"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute("clip-rule", "evenodd"),
        attribute(
          "d",
          "M5 2a3 3 0 0 0-3 3v14a3 3 0 0 0 3 3h14a3 3 0 0 0 3-3V5a3 3 0 0 0-3-3H5Zm1 4a1 1 0 0 0 0 2h5a1 1 0 1 0 0-2H6Zm-1 6a1 1 0 0 1 1-1h12a1 1 0 1 1 0 2H6a1 1 0 0 1-1-1Zm1 4a1 1 0 1 0 0 2h12a1 1 0 1 0 0-2H6Z",
        ),
        attribute("fill-rule", "evenodd"),
      ]),
    ],
  )
}

pub fn speech_bubble(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "currentColor"),
      attribute.role("img"),
      attribute("viewBox", "0 0 24 24"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute(
          "d",
          "M12 22a10 10 0 1 0-8.45-4.64c.13.19.11.44-.04.61l-2.06 2.37A1 1 0 0 0 2.2 22H12Z",
        ),
      ]),
    ],
  )
}

pub fn arrow_out(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "currentColor"),
      attribute("viewBox", "0 0 18 18"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute(
          "d",
          "M9.75 12V4.05751L12.225 6.53251C12.2928 6.6105 12.376 6.67364 12.4693 6.71797C12.5627 6.76229 12.6642 6.78684 12.7675 6.79008C12.8708 6.79332 12.9737 6.77518 13.0696 6.73679C13.1656 6.6984 13.2525 6.6406 13.3251 6.56701C13.3977 6.49341 13.4542 6.40562 13.4912 6.30913C13.5283 6.21265 13.545 6.10955 13.5402 6.00631C13.5355 5.90307 13.5096 5.80192 13.4639 5.7092C13.4183 5.61648 13.3539 5.53419 13.275 5.46751L9.525 1.71751C9.3848 1.58009 9.19631 1.50311 9 1.50311C8.80369 1.50311 8.6152 1.58009 8.475 1.71751L4.725 5.46751C4.60104 5.6101 4.53547 5.79422 4.54139 5.98307C4.54731 6.17192 4.6243 6.35157 4.75695 6.48612C4.8896 6.62066 5.06814 6.70018 5.25688 6.70878C5.44563 6.71738 5.63066 6.65443 5.775 6.53251L8.25 4.05001V12C8.25 12.1989 8.32902 12.3897 8.46967 12.5303C8.61032 12.671 8.80109 12.75 9 12.75C9.19891 12.75 9.38968 12.671 9.53033 12.5303C9.67098 12.3897 9.75 12.1989 9.75 12ZM2.25 15C2.05109 15 1.86032 15.079 1.71967 15.2197C1.57902 15.3603 1.5 15.5511 1.5 15.75C1.5 15.9489 1.57902 16.1397 1.71967 16.2803C1.86032 16.421 2.05109 16.5 2.25 16.5H15.75C15.9489 16.5 16.1397 16.421 16.2803 16.2803C16.421 16.1397 16.5 15.9489 16.5 15.75C16.5 15.5511 16.421 15.3603 16.2803 15.2197C16.1397 15.079 15.9489 15 15.75 15H2.25Z",
        ),
      ]),
    ],
  )
}

pub fn arrow_in(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "currentColor"),
      attribute("viewBox", "0 0 18 18"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute(
          "d",
          "M9.75 2.25V10.1925L12.225 7.7175C12.2928 7.6395 12.376 7.57636 12.4693 7.53204C12.5627 7.48771 12.6642 7.46316 12.7675 7.45992C12.8708 7.45668 12.9737 7.47483 13.0696 7.51321C13.1656 7.5516 13.2525 7.60941 13.3251 7.683C13.3977 7.75659 13.4542 7.84439 13.4912 7.94087C13.5283 8.03736 13.545 8.14046 13.5402 8.24369C13.5355 8.34693 13.5096 8.44809 13.4639 8.5408C13.4183 8.63352 13.3539 8.71581 13.275 8.7825L9.525 12.5325C9.3848 12.6699 9.19631 12.7469 9 12.7469C8.80369 12.7469 8.6152 12.6699 8.475 12.5325L4.725 8.7825C4.60104 8.63991 4.53547 8.45578 4.54139 8.26694C4.54731 8.07809 4.6243 7.89844 4.75695 7.76389C4.8896 7.62935 5.06814 7.54982 5.25688 7.54122C5.44563 7.53262 5.63066 7.59558 5.775 7.7175L8.25 10.2V2.25C8.25 2.05109 8.32902 1.86032 8.46967 1.71967C8.61032 1.57902 8.80109 1.5 9 1.5C9.19891 1.5 9.38968 1.57902 9.53033 1.71967C9.67098 1.86032 9.75 2.05109 9.75 2.25ZM2.25 14.9969C2.05109 14.9969 1.86032 15.0759 1.71967 15.2166C1.57902 15.3572 1.5 15.548 1.5 15.7469C1.5 15.9458 1.57902 16.1366 1.71967 16.2772C1.86032 16.4179 2.05109 16.4969 2.25 16.4969H15.75C15.9489 16.4969 16.1397 16.4179 16.2803 16.2772C16.421 16.1366 16.5 15.9458 16.5 15.7469C16.5 15.548 16.421 15.3572 16.2803 15.2166C16.1397 15.0759 15.9489 14.9969 15.75 14.9969H2.25Z",
        ),
      ]),
    ],
  )
}

pub fn checkmark(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "currentColor"),
      attribute("viewBox", "0 0 18 18"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute("clip-rule", "evenodd"),
        attribute(
          "d",
          "M14.295 5.20499C14.5057 5.41592 14.624 5.70186 14.624 5.99999C14.624 6.29811 14.5057 6.58405 14.295 6.79499L8.29501 12.795C8.08408 13.0057 7.79814 13.124 7.50001 13.124C7.20189 13.124 6.91595 13.0057 6.70501 12.795L3.70501 9.79499C3.50629 9.58172 3.39811 9.29965 3.40325 9.0082C3.40839 8.71675 3.52646 8.43867 3.73258 8.23255C3.9387 8.02643 4.21678 7.90837 4.50823 7.90322C4.79968 7.89808 5.08175 8.00627 5.29501 8.20499L7.50001 10.41L12.705 5.20499C12.916 4.99431 13.2019 4.87598 13.5 4.87598C13.7981 4.87598 14.0841 4.99431 14.295 5.20499Z",
        ),
        attribute("fill-rule", "evenodd"),
      ]),
    ],
  )
}

pub fn message_bubbles(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute.role("img"),
      attribute("viewBox", "0 0 24 24"),
      attribute("fill", "currentColor"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute(
          "d",
          "M18.91 12.98a5.45 5.45 0 0 1 2.18 6.2c-.1.33-.09.68.1.96l.83 1.32a1 1 0 0 1-.84 1.54h-5.5A5.6 5.6 0 0 1 10 17.5a5.6 5.6 0 0 1 5.68-5.5c1.2 0 2.32.36 3.23.98Z",
        ),
      ]),
      svg.path([
        attribute(
          "d",
          "M19.24 10.86c.32.16.72-.02.74-.38L20 10c0-4.42-4.03-8-9-8s-9 3.58-9 8c0 1.5.47 2.91 1.28 4.11.14.21.12.49-.06.67l-1.51 1.51A1 1 0 0 0 2.4 18h5.1a.5.5 0 0 0 .49-.5c0-4.2 3.5-7.5 7.68-7.5 1.28 0 2.5.3 3.56.86Z",
        ),
      ]),
    ],
  )
}

pub fn hashtag(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "currentColor"),
      attribute.role("img"),
      attribute("viewBox", "0 0 24 24"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute("clip-rule", "evenodd"),
        attribute(
          "d",
          "M10.99 3.16A1 1 0 1 0 9 2.84L8.15 8H4a1 1 0 0 0 0 2h3.82l-.67 4H3a1 1 0 1 0 0 2h3.82l-.8 4.84a1 1 0 0 0 1.97.32L8.85 16h4.97l-.8 4.84a1 1 0 0 0 1.97.32l.86-5.16H20a1 1 0 1 0 0-2h-3.82l.67-4H21a1 1 0 1 0 0-2h-3.82l.8-4.84a1 1 0 1 0-1.97-.32L15.15 8h-4.97l.8-4.84ZM14.15 14l.67-4H9.85l-.67 4h4.97Z",
        ),
        attribute("fill-rule", "evenodd"),
      ]),
    ],
  )
}

pub fn chevron_down(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "currentColor"),
      attribute.role("img"),
      attribute("viewBox", "0 0 24 24"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute(
          "d",
          "M5.3 9.3a1 1 0 0 1 1.4 0l5.3 5.29 5.3-5.3a1 1 0 1 1 1.4 1.42l-6 6a1 1 0 0 1-1.4 0l-6-6a1 1 0 0 1 0-1.42Z",
        ),
      ]),
    ],
  )
}

pub fn user_wave(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "currentColor"),
      attribute.role("img"),
      attribute("viewBox", "0 0 24 24"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute("d", "M13 10a4 4 0 1 0 0-8 4 4 0 0 0 0 8Z"),
      ]),
      svg.path([
        attribute(
          "d",
          "M3 5v-.75C3 3.56 3.56 3 4.25 3s1.24.56 1.33 1.25C6.12 8.65 9.46 12 13 12h1a8 8 0 0 1 8 8 2 2 0 0 1-2 2 .21.21 0 0 1-.2-.15 7.65 7.65 0 0 0-1.32-2.3c-.15-.2-.42-.06-.39.17l.25 2c.02.15-.1.28-.25.28H9a2 2 0 0 1-2-2v-2.22c0-1.57-.67-3.05-1.53-4.37A15.85 15.85 0 0 1 3 5Z",
        ),
      ]),
    ],
  )
}

pub fn info_circle(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "currentColor"),
      attribute.role("img"),
      attribute("viewBox", "0 0 24 24"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.circle([
        attribute("cx", "12"),
        attribute("cy", "12"),
        attribute("fill", "transparent"),
        attribute("r", "10"),
      ]),
      svg.path([
        attribute("clip-rule", "evenodd"),
        attribute(
          "d",
          "M23 12a11 11 0 1 1-22 0 11 11 0 0 1 22 0Zm-9.5-4.75a1.25 1.25 0 1 1-2.5 0 1.25 1.25 0 0 1 2.5 0Zm-.77 3.96a1 1 0 1 0-1.96-.42l-1.04 4.86a2.77 2.77 0 0 0 4.31 2.83l.24-.17a1 1 0 1 0-1.16-1.62l-.24.17a.77.77 0 0 1-1.2-.79l1.05-4.86Z",
        ),
        attribute("fill-rule", "evenodd"),
      ]),
    ],
  )
}

pub fn x(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("fill", "currentColor"),
      attribute("viewBox", "0 0 24 24"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute(
          "d",
          "M19.3 20.7a1 1 0 0 0 1.4-1.4L13.42 12l7.3-7.3a1 1 0 0 0-1.42-1.4L12 10.58l-7.3-7.3a1 1 0 0 0-1.4 1.42L10.58 12l-7.3 7.3a1 1 0 1 0 1.42 1.4L12 13.42l7.3 7.3Z",
        ),
      ]),
    ],
  )
}

pub fn discord(attrs: List(attribute.Attribute(a))) -> Element(a) {
  svg.svg(
    [
      attribute("viewBox", "0 0 640 512"),
      attribute("fill", "currentColor"),
      attribute("xmlns", "http://www.w3.org/2000/svg"),
      ..attrs
    ],
    [
      svg.path([
        attribute(
          "d",
          "M524.531,69.836a1.5,1.5,0,0,0-.764-.7A485.065,485.065,0,0,0,404.081,32.03a1.816,1.816,0,0,0-1.923.91,337.461,337.461,0,0,0-14.9,30.6,447.848,447.848,0,0,0-134.426,0,309.541,309.541,0,0,0-15.135-30.6,1.89,1.89,0,0,0-1.924-.91A483.689,483.689,0,0,0,116.085,69.137a1.712,1.712,0,0,0-.788.676C39.068,183.651,18.186,294.69,28.43,404.354a2.016,2.016,0,0,0,.765,1.375A487.666,487.666,0,0,0,176.02,479.918a1.9,1.9,0,0,0,2.063-.676A348.2,348.2,0,0,0,208.12,430.4a1.86,1.86,0,0,0-1.019-2.588,321.173,321.173,0,0,1-45.868-21.853,1.885,1.885,0,0,1-.185-3.126c3.082-2.309,6.166-4.711,9.109-7.137a1.819,1.819,0,0,1,1.9-.256c96.229,43.917,200.41,43.917,295.5,0a1.812,1.812,0,0,1,1.924.233c2.944,2.426,6.027,4.851,9.132,7.16a1.884,1.884,0,0,1-.162,3.126,301.407,301.407,0,0,1-45.89,21.83,1.875,1.875,0,0,0-1,2.611,391.055,391.055,0,0,0,30.014,48.815,1.864,1.864,0,0,0,2.063.7A486.048,486.048,0,0,0,610.7,405.729a1.882,1.882,0,0,0,.765-1.352C623.729,277.594,590.933,167.465,524.531,69.836ZM222.491,337.58c-28.972,0-52.844-26.587-52.844-59.239S193.056,219.1,222.491,219.1c29.665,0,53.306,26.82,52.843,59.239C275.334,310.993,251.924,337.58,222.491,337.58Zm195.38,0c-28.971,0-52.843-26.587-52.843-59.239S388.437,219.1,417.871,219.1c29.667,0,53.307,26.82,52.844,59.239C470.715,310.993,447.538,337.58,417.871,337.58Z",
        ),
      ]),
    ],
  )
}
