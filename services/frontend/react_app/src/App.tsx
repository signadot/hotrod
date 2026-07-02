import './App.css'
import {HomePage} from "./pages/home.tsx";
import {SessionProvider} from "./context/sessionContext/context.tsx";
import {ChakraProvider, extendTheme, type ThemeConfig} from "@chakra-ui/react";

import { accordionAnatomy } from '@chakra-ui/anatomy'
import { createMultiStyleConfigHelpers } from '@chakra-ui/react'

const { definePartsStyle, defineMultiStyleConfig } =
    createMultiStyleConfigHelpers(accordionAnatomy.keys)

const accordionTheme = defineMultiStyleConfig({
    baseStyle: definePartsStyle({
        container: { borderColor: 'gray.600' },
        button: { color: 'whiteAlpha.900', _hover: { bg: 'whiteAlpha.100' } },
        panel: { color: 'whiteAlpha.800' },
    }),
})

const config: ThemeConfig = {
    initialColorMode: 'dark',
    useSystemColorMode: false,
}

const theme = extendTheme({
    config,
    styles: {
        global: {
            body: {
                bg: 'gray.900',
                color: 'whiteAlpha.900',
            },
        },
    },
    components: {
        Accordion: accordionTheme,
        Card: {
            baseStyle: {
                container: {
                    bg: 'gray.800',
                    color: 'whiteAlpha.900',
                    borderColor: 'gray.700',
                },
            },
        },
        Button: {
            variants: {
                outline: {
                    borderColor: 'gray.600',
                    color: 'whiteAlpha.900',
                    _hover: { bg: 'whiteAlpha.100' },
                },
            },
        },
        Select: {
            variants: {
                filled: {
                    field: {
                        bg: 'gray.700',
                        color: 'whiteAlpha.900',
                        _hover: { bg: 'gray.600' },
                        _focusVisible: { bg: 'gray.600' },
                    },
                },
            },
            defaultProps: {
                variant: 'filled',
            },
        },
    },
})

function App() {
    return (
        <ChakraProvider theme={theme}>
            <SessionProvider>
                <HomePage/>
            </SessionProvider>
        </ChakraProvider>
    )
}

export default App
