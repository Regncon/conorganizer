import { AuthProvider } from '@/components/AuthProvider';
import MainNavigator from '@/components/MainNavigator';
import { ThemeProvider } from '@/lib/mui';
import { muiDark } from '@/lib/muiTheme';
import Event from './Event';
type Props = {
    params: { id: string };
};

const page = ({ params }: Props) => {
    const { id } = params;

    return (
        <>
            <Event id={id} />
            <MainNavigator />
        </>
    );
};

export default page;
