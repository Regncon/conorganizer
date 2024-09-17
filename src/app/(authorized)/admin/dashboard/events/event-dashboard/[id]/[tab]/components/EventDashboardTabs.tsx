import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Box from '@mui/material/Box';
import { Hidden, Link, Typography, type SxProps } from '@mui/material';
import GroupIcon from '@mui/icons-material/Group';
import FavoriteIcon from '@mui/icons-material/Favorite';
import Settings from '@mui/icons-material/Settings';
import RoomPreferencesIcon from '@mui/icons-material/RoomPreferences';
import EditIcon from '@mui/icons-material/Edit';
import type { PropsWithChildren } from 'react';
import HideLabel from './ui/HideLabel';

type props = {
    id: string;
    value: number;
};

export default function EventDashboardTabs({ id, value }: props) {
    const tabsSx: SxProps = {
        padding: { md: '1rem', xs: '0' },
    };
    return (
        <Box sx={{ width: '100%' }}>
            <Link href="/admin/dashboard/events/">Tilbake til arrangementer</Link>
            <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                <Tabs
                    sx={tabsSx}
                    value={value}
                    variant="scrollable"
                    scrollButtons
                    allowScrollButtonsMobile
                    allowscrollbuttonsmobilearia-label="Velg side"
                >
                    <Tab
                        icon={<GroupIcon />}
                        sx={tabsSx}
                        iconPosition="start"
                        label={<HideLabel>Spillere</HideLabel>}
                        disabled
                    />
                    <Tab
                        sx={tabsSx}
                        icon={<FavoriteIcon />}
                        iconPosition="start"
                        label={<HideLabel>Ønskeliste</HideLabel>}
                        disabled
                    />
                    <Tab
                        sx={tabsSx}
                        icon={<Settings />}
                        iconPosition="start"
                        label={<HideLabel>Innstillinger</HideLabel>}
                        href={`/admin/dashboard/events/event-dashboard/${id}/settings/`}
                    />
                    <Tab
                        sx={tabsSx}
                        icon={<RoomPreferencesIcon />}
                        iconPosition="start"
                        label={<HideLabel>Rom</HideLabel>}
                        href={`/admin/dashboard/events/event-dashboard/${id}/room/`}
                    />
                    <Tab
                        sx={tabsSx}
                        icon={<EditIcon />}
                        iconPosition="start"
                        label={<HideLabel>Rediger</HideLabel>}
                        href={`/admin/dashboard/events/event-dashboard/${id}/edit/`}
                    />
                </Tabs>
            </Box>
        </Box>
    );
}
