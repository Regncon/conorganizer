'use client';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Box from '@mui/material/Box';

type props = {
    id: string;
    value: number;
};

export default function EventDashboardTabs({ id, value }: props) {
    return (
        <Box sx={{ width: '100%' }}>
            <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                <Tabs value={value} aria-label="Velg side">
                    <Tab label="Spillere" disabled />
                    <Tab label="Ønskeliste" disabled />
                    <Tab label="Innstillinger" href={`/admin/dashboard/events/event-dashboard/${id}/settings/`} />
                    <Tab label="Rom" disabled />
                    <Tab label="Rediger" href={`/admin/dashboard/events/event-dashboard/${id}/edit/`} />
                </Tabs>
            </Box>
        </Box>
    );
}
