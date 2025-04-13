import { Inter, Lora, Source_Sans_3 } from 'next/font/google'

// define your variable fonts
const inter = Inter({ subsets: ['latin'] })
const lora = Lora({subsets: ['latin']})
// define 2 weights of a non-variable font
const sourceCodePro400 = Source_Sans_3({ weight: '400', subsets: ['latin'] })

export { inter, lora, sourceCodePro400 }