import { Paper } from '@mui/material';
import EventDashboardTabs from '../EventDashboardTabs';
import Settings from './Settings';

type Props = {
    params: {
        id: string;
    };
};

const Page = async ({ params: { id } }: Props) => {
    return (
        <Paper>
            <EventDashboardTabs id={id} value={2} />
            <Settings id={id} />
        </Paper>
    );
};

export default Page;
