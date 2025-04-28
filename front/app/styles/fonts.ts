import { Inter, Roboto, Lora, Source_Sans_3, Open_Sans } from 'next/font/google'

// define your variable fonts
const inter = Inter({ weight: "300", subsets: ['latin'] })
const roboto = Roboto({ weight: "300", subsets: ['latin'] })
const lora = Lora({subsets: ['latin']})
const openSans = Open_Sans({
  weight: '300',
  subsets: ['latin']
})
const sourceCodePro400 = Source_Sans_3({
  weight: '300', subsets: ['latin']
})

export { inter, lora, roboto, sourceCodePro400, openSans }