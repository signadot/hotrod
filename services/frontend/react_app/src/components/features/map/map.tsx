import React, {useEffect, useMemo} from "react";
import {MapContainer, TileLayer, useMap} from "react-leaflet";
import "leaflet/dist/leaflet.css";
import * as L from 'leaflet';
import {LatLngTuple} from 'leaflet';
import "leaflet-routing-machine";
import {Flex} from "@chakra-ui/react";

interface RouteMapProps {
    start?: [number, number];
    end?: [number, number];
    center: LatLngTuple;
}

const COORDINATES_MOCKED: Record<number, [number, number]> = {
    1: [37.7749, -122.4194], // My Home (San Francisco, downtown)
    123: [37.7764, -122.4241], // Rachel's Floral Designs (near Hayes Valley)
    392: [37.7723, -122.4108], // Trom Chocolatier (Mission District)
    567: [37.7807, -122.4081], // Amazing Coffee Roasters (SoMa)
    731: [37.7689, -122.4494] // Japanese Desserts (near Golden Gate Park)
};

const getCenter = (start: [number, number] | undefined, end: [number, number] | undefined): LatLngTuple => {
    if (!start || !end) {
        return [37.562304, -122.32668]
    }

    return [(start[0] + end[0]) / 2, (start[1] + end[1]) / 2];
}

const RoutingControl = ({ start, end, center }: RouteMapProps) => {
    const map = useMap();

    useEffect(() => {
        if (!map || !start || !end) return () => {};

        map.setView(center);

        const customPlan = L.Routing.plan([L.latLng(start), L.latLng(end)], {
            createMarker: (i, waypoint) => {
                return L.marker(waypoint.latLng, {
                    icon: L.icon({
                        iconUrl: i === 0
                            ? "https://cdn-icons-png.flaticon.com/512/2991/2991122.png" // Start marker
                            : "https://cdn-icons-png.flaticon.com/512/190/190411.png", // End marker
                        iconSize: [25, 41],
                        iconAnchor: [12, 41],
                    }),
                });
            }, // Prevent default markers
        });

        const routingControl = L.Routing.control({
            waypoints: [L.latLng(start[0], start[1]), L.latLng(end[0], end[1])],
            routeWhileDragging: true,
            addWaypoints: false,
            lineOptions: {
                styles: [{ color: "blue", weight: 6 }], // Wider path
                extendToWaypoints: true, // Default is true, ensures lines extend to waypoints
                missingRouteTolerance: 10, // Default is 10 (meters)
            },
            plan:customPlan,
        }).addTo(map);

        return () => {
            map.removeControl(routingControl);
        };
    }, [map, start, end, center]);

    return null;
};

const RouteMap: React.FC<Omit<RouteMapProps, "center">> = ({start, end}) => {
    const center = useMemo((): LatLngTuple => {
        return getCenter(start, end);
    }, [start, end]);

    return (
        <MapContainer
            center={center}
            zoom={17}
            style={{height: "100%", width: "100%"}}
        >
            <TileLayer
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            />
            <RoutingControl start={start} end={end} center={center}/>
        </MapContainer>
    );
};


type MapProps = {
    dropoffLocationID: number;
    pickupLocationID: number;
}

export const Map = ({dropoffLocationID, pickupLocationID}: MapProps) => {
    const dropoffCoords = COORDINATES_MOCKED[dropoffLocationID as keyof typeof COORDINATES_MOCKED];
    const pickupCoords = COORDINATES_MOCKED[pickupLocationID as keyof typeof COORDINATES_MOCKED];

    return (
        <Flex w='100%' h='100%' borderRadius={16} overflow='hidden'>
            <RouteMap start={pickupCoords} end={dropoffCoords} />
        </Flex>
    );
}