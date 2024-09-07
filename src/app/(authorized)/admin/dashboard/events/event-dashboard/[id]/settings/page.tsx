import { Paper } from '@mui/material';
import EventDashboardTabs from '../components/EventDashboardTabs';
import Settings from './components/Settings';
import { getAllEvents } from '$app/(public)/components/lib/serverAction';

type Props = {
    params: {
        id: string;
    };
};

const Page = async ({ params: { id } }: Props) => {
    const allEvents = await getAllEvents();
    return (
        <Paper sx={{ padding: { sm: '1rem' } }} elevation={0}>
            <EventDashboardTabs id={id} value={2} />
            <Settings id={id} allEvents={allEvents} />
        </Paper>
    );
};

export default Page;
