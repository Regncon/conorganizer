import { Box, Typography } from '@mui/material';
import Image from 'next/image';
import { PoolName, RoomName } from '$lib/enums';
import { getAllEvents } from '$app/(public)/components/lib/serverAction';
import RoomMapItem from './RoomMapItem';
import { RoomItemInfo } from '$lib/types';

type Props = {
    pool: PoolName | undefined;
};

const RoomMap = async ({ pool }: Props) => {
    if (!pool) {
        return <Typography>Velg en pulje</Typography>;
    }
    const events = await getAllEvents();

    const poolTitles = {
        [PoolName.fridayEvening]: 'Fredag Kveld',
        [PoolName.saturdayMorning]: 'Lørdag Morgen',
        [PoolName.saturdayEvening]: 'Lørdag Kveld',
        [PoolName.sundayMorning]: 'Søndag Morgen',
    };

    /*     <RoomAddButton
                events={events}
                roomCoordinates={{ x: 2560, y: 450 }}
                poolName={poolName}
                roomName={RoomName.Klang}
            />
            <RoomAddButton
                events={events}
                roomCoordinates={{
                    x: 2532,
                    y: 628,
                }}
                poolName={poolName}
                roomName={RoomName.Sonate}
            />
            <RoomAddButton
                events={events}
                roomCoordinates={{
                    x: 2510,
                    y: 800,
                }}
                poolName={poolName}
                roomName={RoomName.Ballade}
            />
            <RoomAddButton
                events={events}
                roomCoordinates={{
                    x: 2490,
                    y: 980,
                }}
                poolName={poolName}
                roomName={RoomName.Klaver}
            />
            <RoomAddButton
                events={events}
                roomCoordinates={{
                    x: 2650,
                    y: 1125,
                }}
                poolName={poolName}
                roomName={RoomName.Hymne}
            />
            <RoomAddButton
                events={events}
                roomCoordinates={{
                    x: 2650,
                    y: 1280,
                }}
                poolName={poolName}
                roomName={RoomName.Fanfare}
            />
            <RoomAddButton
                events={events}
                roomCoordinates={{
                    x: 2460,
                    y: 1510,
                }}
                poolName={poolName}
                roomName={RoomName.Kammer}
            />
            <RoomAddButton
                events={events}
                roomCoordinates={{
                    x: 2097,
                    y: 1950,
                }}
                poolName={poolName}
                roomName={RoomName.Beyer}
            />
            <RoomAddButton
                events={events}
                roomCoordinates={{
                    x: 1900,
                    y: 2050,
                }}
                poolName={poolName}
                roomName={RoomName.Siljuslåtten}
            />
            <RoomAddButton
                events={events}
                roomCoordinates={{
                    x: 1560,
                    y: 2050,
                }}
                poolName={poolName}
                roomName={RoomName.PeerGynt}
            />
            <RoomAddButton
                events={events}
                roomCoordinates={{
                    x: 250,
                    y: 1000,
                }}
                poolName={poolName}
                roomName={RoomName.SolveigsSang}
            />
            <RoomAddButton
                events={events}
                roomCoordinates={{
                    x: 250,
                    y: 750,
                }}
                poolName={poolName}
                roomName={RoomName.AnitrasDans}
            />
            <RoomAddButton
                events={events}
                roomCoordinates={{
                    x: 450,
                    y: 550,
                }}
                poolName={poolName}
                roomName={RoomName.Werenskiold}
            />
            <RoomAddButton
                events={events}
                roomCoordinates={{
                    x: 450,
                    y: 250,
                }}
                poolName={poolName}
                roomName={RoomName.Welhaven}
            />
    */

    const mapItems: RoomItemInfo[] = [
        { roomName: RoomName.NotSet, top: 350, left: 900 },
        { roomName: RoomName.Styrerom1, top: 350, left: 1500 },
        { roomName: RoomName.Styrerom2, top: 550, left: 1500 },
        { roomName: RoomName.Styrerom3, top: 750, left: 1500 },
        { roomName: RoomName.Styrerom4, top: 350, left: 1800 },
        { roomName: RoomName.Styrerom5, top: 550, left: 1800 },
        { roomName: RoomName.Styrerom6, top: 750, left: 1800 },
        { roomName: RoomName.Klang, top: 450, left: 2460 },
        { roomName: RoomName.Sonate, top: 640, left: 2460 },
        { roomName: RoomName.Ballade, top: 800, left: 2460 },
        { roomName: RoomName.Klaver, top: 980, left: 2460 },
        { roomName: RoomName.Hymne, top: 1125, left: 2650 },
        { roomName: RoomName.Fanfare, top: 1280, left: 2650 },
        { roomName: RoomName.Kammer, top: 1510, left: 2460 },
        { roomName: RoomName.Beyer, top: 1950, left: 2097 },
        { roomName: RoomName.Siljuslåtten, top: 2050, left: 1900 },
        { roomName: RoomName.PeerGynt, top: 2050, left: 1560 },
        { roomName: RoomName.SolveigsSang, top: 1000, left: 250 },
        { roomName: RoomName.AnitrasDans, top: 750, left: 250 },
        { roomName: RoomName.Werenskiold, top: 550, left: 450 },
        { roomName: RoomName.Welhaven, top: 250, left: 450 },
    ];

    return (
        <Box>
            <Typography
                variant="h1"
                sx={{ fontSize: '90px', color: 'black', position: 'absolute', top: '100px', left: '900px' }}
            >
                {poolTitles[pool]}
            </Typography>
            <Typography variant="h2" sx={{ color: 'black', position: 'absolute', top: '300px', left: '900px' }}>
                Arrangementer på {poolTitles[pool]} uten rom
            </Typography>
            {mapItems.map((roomItem) => {
                return (
                    <RoomMapItem
                        roomName={roomItem.roomName}
                        key={roomItem.roomName}
                        top={roomItem.top}
                        left={roomItem.left}
                        poolName={pool}
                        events={events}
                    ></RoomMapItem>
                );
            })}
            <Image src={'/rooms.webp'} alt={'Romkart'} width={'2901'} height={'2073'}></Image>
        </Box>
    );
};
export default RoomMap;
