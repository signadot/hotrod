import {ReactNode} from "react";
import {Stack} from "@chakra-ui/react";
import {Header} from "./common/header.tsx";


type MainLayoutProps = {
    titleSuffix: string,
    children: ReactNode,
}

export const MainLayout = ({ titleSuffix, children }: MainLayoutProps) => {
    return (
        <Stack h='100vh' w='100vw' px={12} py={8}>
            <Header titleSuffix={titleSuffix} />
            <Stack mt={12} h='100%'>
                {children}
            </Stack>
        </Stack>
    )
}