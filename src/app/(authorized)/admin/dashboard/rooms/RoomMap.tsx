import { Box, Typography } from '@mui/material';
import Image from 'next/image';
import RoomCard from './RoomCard';
import { PoolName, RoomName } from '$lib/enums';
import { getAllEvents } from '$app/(public)/components/lib/serverAction';

type Props = {
    pool: PoolName;
};

const RoomMap = async ({ pool }: Props) => {
    const events = await getAllEvents();

    const poolTitles = {
        [PoolName.fridayEvening]: 'Fredag Kveld',
        [PoolName.saturdayMorning]: 'Lørdag Morgen',
        [PoolName.saturdayEvening]: 'Lørdag Kveld',
        [PoolName.sundayMorning]: 'Søndag Morgen',
    };

    return (
        <Box>
            <Typography
                variant="h1"
                sx={{ fontSize: '90px', color: 'black', position: 'absolute', top: '100px', left: '900px' }}
            >
                {poolTitles[pool]}
            </Typography>
            <Typography variant="h2" sx={{ color: 'black', position: 'absolute', top: '300px', left: '1000px' }}>
                Styrerom 1 Gerhard
            </Typography>
            <RoomCard
                poolName={pool}
                roomName={RoomName.Styreromm1}
                title={'Kjempegøy drager og fangehull'}
                gameMaster={'Kari Nordmann'}
                system={'D&D'}
                imageUri="/blekksprut2.jpg"
                events={events}
            ></RoomCard>
            <RoomCard
                poolName={pool}
                roomName={RoomName.Klang}
                title={'Kjempegøy drager og fangehull'}
                gameMaster={'Kari Nordmann'}
                system={'D&D'}
                imageUri="/blekksprut2.jpg"
                events={events}
            ></RoomCard>
            <RoomCard
                poolName={pool}
                roomName={RoomName.Sonate}
                title={'En telefon fra Cthulhu'}
                gameMaster={'Ola Nordmann'}
                system={'Call of Cthulhu'}
                imageUri="/blekksprut2.jpg"
                events={events}
            ></RoomCard>
            <Image src={'/rooms.webp'} alt={'Romkart'} width={'2901'} height={'2073'}></Image>
        </Box>
    );
};
export default RoomMap;
