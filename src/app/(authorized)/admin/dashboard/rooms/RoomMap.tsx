import { Box, Typography } from '@mui/material';
import Image from 'next/image';
import { PoolName, RoomName } from '$lib/enums';
import { getAllEvents } from '$app/(public)/components/lib/serverAction';
import RoomMapItem from './RoomMapItem';
import { ConEvent } from '$lib/types';

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

    return (
        <Box>
            <Typography
                variant="h1"
                sx={{ fontSize: '90px', color: 'black', position: 'absolute', top: '100px', left: '900px' }}
            >
                {poolTitles[pool]}
            </Typography>
            <Typography variant="h2" sx={{ color: 'black', position: 'absolute', top: '300px', left: '1000px' }}>
                Arrangementer på {poolTitles[pool]} uten rom
            </Typography>
            {events.map((event) => {
                return (
                    <RoomMapItem key={event.id} eventId={event.id || ''} poolName={pool} events={events}></RoomMapItem>
                );
            })}
            <Image src={'/rooms.webp'} alt={'Romkart'} width={'2901'} height={'2073'}></Image>
        </Box>
    );
};
export default RoomMap;
