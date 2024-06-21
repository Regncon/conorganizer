import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import { useFormStatus } from 'react-dom';
import { disableAndLoadingSpinner } from './LoginPage';

type Props = {
    disabled?: boolean;
    endIcon?: React.ReactNode;
};

const LoginButton = ({}: Props) => {
    const { pending } = useFormStatus();

    return (
        <Button type="submit" {...disableAndLoadingSpinner(true, pending)}>
            Logg inn
        </Button>
    );
};

export default LoginButton;
