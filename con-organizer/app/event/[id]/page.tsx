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
            <Event id={id} />
        </ThemeProvider>
    );
};

export default page;
