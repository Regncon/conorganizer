import { Paper } from '@mui/material';
import EventDashboardTabs from '../EventDashboardTabs';

type Props = {
    params: {
        id: string;
    };
};

const Page = async ({ params: { id } }: Props) => {
    return (
        <Paper>
            <EventDashboardTabs id={id} value={2} />
        </Paper>
    );
};

export default Page;
