import {ReactNode} from "react";
import {Stack} from "@chakra-ui/react";
import {Header} from "./common/header.tsx";


type MainLayoutProps = {
    children: ReactNode,
}

export const MainLayout = ({ children }: MainLayoutProps) => {
    return (
        <Stack h='100vh' w='100vw' px={12} py={8}>
            <Header />
            <Stack mt={12} h='100%'>
                {children}
            </Stack>
        </Stack>
    )
}