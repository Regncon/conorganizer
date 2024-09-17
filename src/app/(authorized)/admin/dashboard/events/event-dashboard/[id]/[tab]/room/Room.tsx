import { getEventById } from '$app/(public)/components/lib/serverAction';
import { PoolName } from '$lib/enums';
import { Box, Link, Paper, Typography } from '@mui/material';

type Props = {
    id: string;
};

const Room = async ({ id }: Props) => {
    const events = await getEventById(id);

    const roomsGrouped = Object.groupBy(events.roomIds, (room) => room.poolName);

    console.log(Object.entries(roomsGrouped));
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
            <Link component={Link} href="/admin/dashboard/rooms">
                Gå til puljetildeling
            </Link>
            {events.poolIds.length === 0 ?
                <Typography variant="h1">Ingen valgt rom</Typography>
            :   null}
            {Object.entries(roomsGrouped)
                .sort()
                .map(([day, rooms]) => {
                    console.log(rooms);

                    return (
                        <Box key={day}>
                            <Typography variant="h1">{translateDay(day as PoolName)}</Typography>
                            {rooms.map((room) => {
                                return (
                                    <Paper
                                        key={room.id}
                                        sx={{ marginBlockEnd: '1rem', marginInline: '1rem', paddingInline: '1rem' }}
                                    >
                                        <Typography variant="h3">{room.roomName}</Typography>
                                        <Typography marginBlockEnd={1}>{room.poolName}</Typography>
                                        <Typography>{room.id}</Typography>
                                        <Typography marginBlockEnd={1}>{room.poolId}</Typography>
                                        <Typography>{room.updateAt}</Typography>
                                        <Typography>{room.updatedBy}</Typography>
                                        <Typography>{room.createdAt}</Typography>
                                        <Typography>{room.createdBy}</Typography>
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
