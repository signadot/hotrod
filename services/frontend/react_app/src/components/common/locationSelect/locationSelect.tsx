import {Select} from "@chakra-ui/react";
import {Location} from "../../../types/location.ts";


type LocationSelectProps = {
    placeholder: string;
    locations: Location[];
    onSelect: (locationID: number) => void;
    selectedLocationID: number | undefined;
}

export const LocationSelect = ({ placeholder, locations, selectedLocationID, onSelect }: LocationSelectProps) => {

    const handleSelect = (event: React.ChangeEvent<HTMLSelectElement>) => {
        onSelect(parseInt(event.target.value))
    }

    return (
        <Select
            placeholder={placeholder}
            variant='filled'
            value={selectedLocationID}
            onChange={handleSelect}
        >
            { locations.map(loc => {
                return (
                    <option value={loc.id}>{loc.name}</option>
                )
            })}
        </Select>
    )
}