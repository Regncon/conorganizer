import { AuthProvider } from '@/components/AuthProvider';
import MainNavigator from '@/components/mainNavigator';
import { ThemeProvider } from '@/lib/mui';
import { muiDark } from '@/lib/muiTheme';
import Event from './Event';
type Props = {
    params: { id: string };
};

const page = ({ params }: Props) => {
    const { id } = params;

    return (
        <ThemeProvider theme={muiDark}>
            <AuthProvider>
                <Event id={id} />
                <MainNavigator />
            </AuthProvider>
        </ThemeProvider>
    );
};

export default page;
