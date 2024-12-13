import './App.css'
import {HomePage} from "./pages/home.tsx";
import {SessionProvider} from "./context/sessionContext/context.tsx";
import {ChakraProvider, extendTheme} from "@chakra-ui/react";

function App() {

    const components = {
        Drawer: {
            variants: {
                aside: {
                    overlay: {
                        pointerEvents: 'none',
                        background: 'transparent',
                    },
                    dialogContainer: {
                        pointerEvents: 'none',
                        background: 'transparent',
                    },
                    dialog: {
                        pointerEvents: 'auto',
                    },
                },
            },
        },
    };

    const theme = extendTheme({components});

    return (
        <ChakraProvider theme={theme}>
            <SessionProvider>
                <HomePage/>
            </SessionProvider>
        </ChakraProvider>
    )
}

export default App
