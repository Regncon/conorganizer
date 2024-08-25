import MainAppBarDesktop from './mainAppBarDesktop';
import MainAppBarMobile from './mainAppBarMobile';

type Props = {};

const MainAppBar = ({ }: Props) => {
    return (
        <>
            <MainAppBarDesktop />
            <MainAppBarMobile />
        </>
    );
};

export default MainAppBar;
