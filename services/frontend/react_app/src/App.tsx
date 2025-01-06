import './App.css'
import {HomePage} from "./pages/home.tsx";
import {SessionProvider} from "./context/sessionContext/context.tsx";
import {ChakraProvider, extendTheme} from "@chakra-ui/react";


import { accordionAnatomy } from '@chakra-ui/anatomy'
import { createMultiStyleConfigHelpers } from '@chakra-ui/react'

const { definePartsStyle, defineMultiStyleConfig } =
    createMultiStyleConfigHelpers(accordionAnatomy.keys)

const baseStyle = definePartsStyle({
    container: {
        borderColor: 'gray.400',
    },
})


const accordionTheme = defineMultiStyleConfig({ baseStyle })

function App() {

    const theme = extendTheme({
        components: { Accordion: accordionTheme },
    })

    return (
        <ChakraProvider theme={theme}>
            <SessionProvider>
                <HomePage/>
            </SessionProvider>
        </ChakraProvider>
    )
}

export default App
