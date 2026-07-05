import { Roboto, Source_Sans_3, Open_Sans } from 'next/font/google'

const roboto = Roboto({ weight: ["300", "700"], subsets: ['latin'] })
const openSans = Open_Sans({
  weight: ["300", "700"],
  subsets: ['latin']
})
/** Source Sans 3 — previously misnamed `sourceCodePro400` despite not being a monospace font. */
const sourceSans400 = Source_Sans_3({
  weight: '300', subsets: ['latin']
})

export { roboto, sourceSans400, openSans }