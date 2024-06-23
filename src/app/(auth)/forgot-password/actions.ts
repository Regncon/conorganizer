import { z } from 'zod';
import type { InitialForgotFormState } from './ForgotPassword';

export const validateForgotFormAction = async (formData: FormData): Promise<InitialForgotFormState> => {
    const formDataEntries = Object.fromEntries(formData);

    const schemaEmail = z.string().email({ message: 'Ugyldig e-post' });

    const resultEmail = schemaEmail.safeParse(formDataEntries.email);
    const passwordError = resultEmail.error?.format()._errors[0];

    const resetErrors: InitialForgotFormState = {
        emailError: '',
    };

    if (!resultEmail.success) {
        return {
            emailError: passwordError ?? resetErrors.emailError,
        };
    }

    return resetErrors;
};
