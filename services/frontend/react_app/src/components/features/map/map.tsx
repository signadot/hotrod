import {Flex} from "@chakra-ui/react";
import {Marker, Map as PigeonMap} from "pigeon-maps";


export const Map = () => {
    return (
        <Flex w='100%' h='100%' borderRadius={16} overflow='hidden'>
            <PigeonMap  defaultCenter={[37.562304, -122.32668]} defaultZoom={17}>
                <Marker width={50} anchor={[37.562304, -122.32668]} />
            </PigeonMap>
        </Flex>
    )
}