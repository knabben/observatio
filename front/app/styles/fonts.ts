import { Inter, Lora, Source_Sans_3, Open_Sans } from 'next/font/google'

// define your variable fonts
const inter = Inter({ subsets: ['latin'] })
const lora = Lora({subsets: ['latin']})
const openSans = Open_Sans({
  weight: '300',
  subsets: ['latin']
})
const sourceCodePro400 = Source_Sans_3({
  weight: '500', subsets: ['latin']
})

export { inter, lora, sourceCodePro400, openSans }