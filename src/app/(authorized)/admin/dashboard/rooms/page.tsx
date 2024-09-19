import { AppBar, Paper, Tab, Tabs, Toolbar, Typography } from '@mui/material';
import { PoolName } from '$lib/enums';
import { redirect } from 'next/navigation';
import RoomMap from './components/RoomMap';
import type { Metadata } from 'next';
import { translatedDays } from '$app/(public)/components/lib/helpers/translation';

type Props = { searchParams: { pool: PoolName | undefined } };
export async function generateMetadata({ searchParams: { pool } }: Props): Promise<Metadata> {
    if (pool) {
        return {
            title: `Romfordeling for ${translatedDays.get(pool)}`,
        };
    }
    return {};
}
const Rooms = async ({ searchParams }: Props) => {
    if (!searchParams.pool) {
        const detfaultPoolPage = `./rooms?pool=${PoolName[PoolName.fridayEvening]}`;
        redirect(detfaultPoolPage);
    }

    let value = 0;
    if (searchParams.pool === PoolName[PoolName.fridayEvening]) {
        value = 0;
    } else if (searchParams.pool === PoolName[PoolName.saturdayMorning]) {
        value = 1;
    } else if (searchParams.pool === PoolName[PoolName.saturdayEvening]) {
        value = 2;
    } else if (searchParams.pool === PoolName[PoolName.sundayMorning]) {
        value = 3;
    }

    return (
        <Paper
            sx={{
                width: 'calc(2901px + 300px + 2rem)',
                height: 'calc(2073px + 300px + 7rem)',
                position: 'absolute',
                left: '0',
                top: '60px',
                padding: '1rem',
                margin: '1rem',
                backgroundColor: 'white',
            }}
        >
            <AppBar position="fixed" sx={{ paddingTop: '60px' }}>
                <Toolbar>
                    <Typography variant="h1">Romfordeling </Typography>
                    <Tabs value={value} aria-label="basic tabs example">
                        <Tab label="Fredag Kveld" href={`./rooms?pool=${PoolName[PoolName.fridayEvening]}`} />
                        <Tab label="Lørdag Morgen" href={`./rooms?pool=${PoolName[PoolName.saturdayMorning]}`} />
                        <Tab label="Lørdag Kveld" href={`./rooms?pool=${PoolName[PoolName.saturdayEvening]}`} />
                        <Tab label="Søndag Morgen" href={`./rooms?pool=${PoolName[PoolName.sundayMorning]}`} />
                    </Tabs>
                </Toolbar>
            </AppBar>
            <Toolbar />
            <RoomMap pool={searchParams.pool as PoolName} />
        </Paper>
    );
};

export default Rooms;
