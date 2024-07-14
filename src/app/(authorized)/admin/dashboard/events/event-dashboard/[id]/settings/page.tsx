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
            <EventDashboardTabs id={id} />
        </Paper>
    );
};

export default Page;
