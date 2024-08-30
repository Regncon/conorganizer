import MainAppBarContent from './MainAppBarContent';

type Props = {
    user: boolean;
    admin: boolean;
};

const MainAppBarMobile = ({ user, admin }: Props) => {
    return (
        <MainAppBarContent
            user={!!user}
            mobile={true}
            admin={admin}
            sx={{ top: 'auto', bottom: 0, display: { xs: 'flex', sm: 'none' } }}
        />
    );
};

export default MainAppBarMobile;
