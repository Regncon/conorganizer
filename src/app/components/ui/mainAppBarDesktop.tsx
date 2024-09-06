import MainAppBarContent from './MainAppBarContent';

type Props = {
    isLoggedIn: boolean;
    admin: boolean;
};

const MainAppBarDesktop = ({ isLoggedIn, admin }: Props) => {
    return (
        <MainAppBarContent
            isLoggedIn={isLoggedIn}
            mobile={false}
            admin={admin}
            sx={{ display: { xs: 'none', sm: 'flex' } }}
        />
    );
};

export default MainAppBarDesktop;
