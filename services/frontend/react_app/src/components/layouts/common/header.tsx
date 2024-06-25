import { Flex, Heading, HStack, Image,} from "@chakra-ui/react";


type HeaderProps = {
    titleSuffix: string,
}

export const Header = ({ titleSuffix }: HeaderProps) => {
    return (
        <Flex w='100%' borderBottom={12}>
            <HStack>
                <Image src='/web_assets/hotrod_logo.png' h={20} w={20}/>
                <Heading>Hotrod Demo App {titleSuffix}</Heading>
                <Heading as='h6' size='xs' justifySelf='self-end' placeSelf='flex-end'>by Signadot</Heading>
            </HStack>
        </Flex>
    )
}