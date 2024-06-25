import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import { useFormStatus } from 'react-dom';
import { disableAndLoadingSpinner } from './LoginPage';

type Props = {
    disabled?: boolean;
};

const LoginButton = ({ disabled = false }: Props) => {
    const { pending } = useFormStatus();

    return (
        <Button type="submit" {...disableAndLoadingSpinner(pending, pending || disabled)}>
            Logg inn
        </Button>
    );
};

export default LoginButton;
