import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Box from '@mui/material/Box';
import { Hidden, Link } from '@mui/material';
import GroupIcon from '@mui/icons-material/Group';
import FavoriteIcon from '@mui/icons-material/Favorite';
import Settings from '@mui/icons-material/Settings';
import RoomPreferencesIcon from '@mui/icons-material/RoomPreferences';
import EditIcon from '@mui/icons-material/Edit';

type props = {
    id: string;
    value: number;
};

export default function EventDashboardTabs({ id, value }: props) {
    return (
        <Box sx={{ width: '100%' }}>
            <Link href="/admin/dashboard/events/">Tilbake til arrangementer</Link>
            <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                <Tabs
                    value={value}
                    variant="scrollable"
                    scrollButtons
                    allowScrollButtonsMobile
                    allowScrollButtonsMobilearia-label="Velg side"
                >
                    <Tab icon={<GroupIcon />} iconPosition="start" label={<Hidden xsUp>Spillere</Hidden>} disabled />
                    <Tab
                        icon={<FavoriteIcon />}
                        iconPosition="start"
                        label={<Hidden xsUp>Ønskeliste</Hidden>}
                        disabled
                    />
                    <Tab
                        icon={<Settings />}
                        iconPosition="start"
                        label={<Hidden xsUp>Innstillinger</Hidden>}
                        href={`/admin/dashboard/events/event-dashboard/${id}/settings/`}
                    />
                    <Tab
                        icon={<RoomPreferencesIcon />}
                        iconPosition="start"
                        label={<Hidden xsUp>Rom</Hidden>}
                        disabled
                    />
                    <Tab
                        icon={<EditIcon />}
                        iconPosition="start"
                        label={<Hidden xsUp>Rediger</Hidden>}
                        href={`/admin/dashboard/events/event-dashboard/${id}/edit/`}
                    />
                </Tabs>
            </Box>
        </Box>
    );
}
