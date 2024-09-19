import { Box, Typography } from '@mui/material';
import Image from 'next/image';
import { PoolName, RoomName } from '$lib/enums';
import { getAllEvents } from '$app/(public)/components/lib/serverAction';
import { RoomItemInfo } from '$lib/types';
import RoomMapItem from './components/RoomMapItem';
import { poolTitles } from './lib/helpers';

type Props = {
    pool: PoolName | undefined;
};

const RoomMap = async ({ pool }: Props) => {
    if (!pool) {
        return <Typography>Velg en pulje</Typography>;
    }
    const events = await getAllEvents();
    const smallRoomsLeft = 2550;
    const styreromLeft = 100;
    const styresomTopStart = 1400;

    const mapItems: RoomItemInfo[] = [
        { roomName: RoomName.NotSet, top: 350, left: 900 },
        { roomName: RoomName.Styrerom1, top: styresomTopStart, left: styreromLeft },
        { roomName: RoomName.Styrerom2, top: styresomTopStart + 200, left: styreromLeft },
        { roomName: RoomName.Styrerom3, top: styresomTopStart + 400, left: styreromLeft },
        { roomName: RoomName.Styrerom4, top: styresomTopStart, left: styreromLeft + 600 },
        { roomName: RoomName.Styrerom5, top: styresomTopStart + 200, left: styreromLeft + 600 },
        { roomName: RoomName.Styrerom6, top: styresomTopStart + 400, left: styreromLeft + 600 },
        { roomName: RoomName.fellerom1, top: 1100, left: 1600 },
        { roomName: RoomName.fellerom2, top: 1450, left: 1900 },
        { roomName: RoomName.Klang, top: 300, left: smallRoomsLeft },
        { roomName: RoomName.Sonate, top: 500, left: smallRoomsLeft },
        { roomName: RoomName.Ballade, top: 700, left: smallRoomsLeft },
        { roomName: RoomName.Klaver, top: 900, left: smallRoomsLeft },
        { roomName: RoomName.Hymne, top: 1100, left: smallRoomsLeft },
        { roomName: RoomName.Fanfare, top: 1300, left: smallRoomsLeft },
        { roomName: RoomName.Kammer, top: 1500, left: smallRoomsLeft },
        { roomName: RoomName.Beyer, top: 1800, left: 2400 },
        { roomName: RoomName.Siljuslåtten, top: 2050, left: 1800 },
        { roomName: RoomName.PeerGynt, top: 2050, left: 1200 },
        { roomName: RoomName.Dovregubben, top: 500, left: 10 },
        { roomName: RoomName.SolveigsSang, top: 1000, left: 100 },
        { roomName: RoomName.AnitrasDans, top: 750, left: 100 },
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
                    />
                );
            })}
            <Image src={'/rooms.webp'} alt={'Romkart'} width={'2901'} height={'2073'}></Image>
        </Box>
    );
};
export default RoomMap;
