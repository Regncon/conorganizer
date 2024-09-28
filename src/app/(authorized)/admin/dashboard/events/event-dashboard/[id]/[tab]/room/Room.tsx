import { getEventById } from '$app/(public)/components/lib/serverAction';
import { PoolName } from '$lib/enums';
import type { RoomChildRef } from '$lib/types';
import { Box, Link, Paper, Typography } from '@mui/material';
import NextLink from 'next/link';
type Props = {
    id: string;
};

const Room = async ({ id }: Props) => {
    const event = await getEventById(id);
    const poolRoomsSorted = new Map<PoolName, RoomChildRef[]>([
        [PoolName.fridayEvening, []],
        [PoolName.saturdayMorning, []],
        [PoolName.saturdayEvening, []],
        [PoolName.sundayMorning, []],
    ]);
    event.roomIds.forEach((room) => {
        const currentPoolRooms = poolRoomsSorted.get(room.poolName);

        if (currentPoolRooms !== undefined) {
            poolRoomsSorted.set(room.poolName, [...currentPoolRooms, room]);
        }
    });

    const translateDay = (day: PoolName) => {
        switch (day) {
            case PoolName.fridayEvening:
                return 'Fredag kveld';
            case PoolName.saturdayMorning:
                return 'Lørdag morgen';
            case PoolName.saturdayEvening:
                return 'Lørdag kveld';
            case PoolName.sundayMorning:
                return 'Søndag morgen';
            default:
                return 'Ukendt';
        }
    };
    return (
        <Box>
            <Link component={NextLink} href="/admin/dashboard/rooms">
                Gå til puljetildeling
            </Link>
            {event.roomIds.length === 0 ?
                <Typography variant="h1">Ingen valgt rom</Typography>
            :   null}
            {[...poolRoomsSorted.entries()].map(([day, rooms]) => {
                return (
                    <Box key={day} sx={{ marginBlock: '2rem' }}>
                        <Typography variant="h1">{translateDay(day)}</Typography>
                        {rooms.map((room) => {
                            return (
                                <Paper
                                    key={room.id}
                                    sx={{ marginBlockEnd: '1rem', marginInline: '1rem', paddingInline: '1rem' }}
                                >
                                    <Typography variant="h3">{`Rom: ${room.roomName}`}</Typography>
                                </Paper>
                            );
                        })}
                    </Box>
                );
            })}
        </Box>
    );
};

export default Room;
