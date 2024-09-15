import { Paper } from '@mui/material';
import EventDashboardTabs from '../components/EventDashboardTabs';
import Settings from './components/Settings';

type Props = {
    params: {
        id: string;
    };
};

const Page = async ({ params: { id } }: Props) => {
    return (
        <Paper sx={{ padding: { sm: '1rem' } }} elevation={0}>
            <EventDashboardTabs id={id} value={2} />
            <Settings id={id} />
        </Paper>
    );
};

export default Page;
