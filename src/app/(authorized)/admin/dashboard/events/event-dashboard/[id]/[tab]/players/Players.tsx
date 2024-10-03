import { Box, Paper, Tab, Tabs, Typography } from '@mui/material';
import { getEventById } from '$app/(public)/components/lib/serverAction';
import PoolNameTabs from './components/PoolNameTabs';
import PlayerManagement from './components/PlayerManagement';
import type { PoolName } from '$lib/enums';

type Props = {
    id: string;
    activeTab: PoolName;
};

const Players = async ({ id, activeTab }: Props) => {
    const event = await getEventById(id);
    const poolId = event.poolIds.find((pool) => pool.poolName === activeTab)?.id;

    return (
        <Box>
            <Typography variant="h1">Spillere:</Typography>
            <Paper sx={{ padding: '1rem' }}>
                <PoolNameTabs id={id} />
                <PlayerManagement id={poolId} poolName={activeTab} />
            </Paper>
        </Box>
    );
};

export default Players;
