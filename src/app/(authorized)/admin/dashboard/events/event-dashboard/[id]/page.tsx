import { Paper } from '@mui/material';
import EventDashboardTabs from './components/EventDashboardTabs';

type Props = {
    params: {
        id: string;
    };
};

const Page = async ({ params: { id } }: Props) => {
    return (
        <Paper>
            <EventDashboardTabs id={id} value={4} />
        </Paper>
    );
};

export default Page;
