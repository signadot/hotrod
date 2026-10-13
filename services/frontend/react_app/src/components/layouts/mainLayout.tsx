import {ReactNode} from "react";
import {Box, Flex} from "@chakra-ui/react";
import {Header} from "./common/header.tsx";

type MainLayoutProps = {
    children: ReactNode,
}

export const MainLayout = ({ children }: MainLayoutProps) => {
    return (
        <Flex direction='column' h='100vh' w='100vw' bg='gray.900'>
            <Header />
            <Box flex={1} overflow='hidden'>
                {children}
            </Box>
        </Flex>
    )
}