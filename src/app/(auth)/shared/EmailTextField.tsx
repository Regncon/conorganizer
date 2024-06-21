import { TextField, InputAdornment } from '@mui/material';
import { emailRegExp } from './utils';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
type Props = {
    defaultValue?: string;
    error?: string;
    helperText?: string;
};

const EmailTextField = ({ defaultValue, error, helperText }: Props) => {
    return (
        <TextField
            type="email"
            name="email"
            autoComplete="email"
            label="e-post"
            variant="outlined"
            error={!!error}
            helperText={helperText}
            defaultValue={defaultValue}
            fullWidth
            required
            InputProps={{
                endAdornment: (
                    <InputAdornment position="end">
                        <AccountCircleIcon />
                    </InputAdornment>
                ),
            }}
            // inputProps={{
            //     pattern: emailRegExp.source,
            //     title: 'epost@example.com',
            // }}
        />
    );
};

export default EmailTextField;
