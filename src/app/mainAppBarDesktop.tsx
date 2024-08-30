import MainAppBarContent from './MainAppBarContent';

type Props = {
    user: boolean;
    admin: boolean;
};

const MainAppBarDesktop = ({ user, admin }: Props) => {
    return (
        <MainAppBarContent user={!!user} mobile={false} admin={admin} sx={{ display: { xs: 'none', sm: 'flex' } }} />
    );
};

export default MainAppBarDesktop;
