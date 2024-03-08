import { Flex, Heading, HStack, Image,} from "@chakra-ui/react";


export const Header = () => {
    return (
        <Flex w='100%' borderBottom={12}>
            <HStack>
                <Image src='/web_assets/hotrod_logo.png' h={20} w={20}/>
                <Heading>Hotrod Demo App</Heading>
                <Heading as='h6' size='xs' justifySelf='self-end' placeSelf='flex-end'>by Signadot</Heading>
            </HStack>
        </Flex>
    )
}