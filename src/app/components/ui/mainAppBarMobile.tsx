import MainAppBarContent from './MainAppBarContent';

type Props = {
    isLoggedIn: boolean;
    admin: boolean;
};

const MainAppBarMobile = ({ isLoggedIn, admin }: Props) => {
    return (
        <MainAppBarContent
            isLoggedIn={isLoggedIn}
            mobile={true}
            admin={admin}
            sx={{
                top: 'auto',
                bottom: 0,
                display: { xs: 'flex', sm: 'none' },
                minHeight: '5.3571rem',
                placeContent: 'center',
            }}
        />
    );
};

export default MainAppBarMobile;
