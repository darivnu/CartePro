import { use, useState } from 'react'

export function PartnerRegisterPage() {
    const [businessName, setBusinessName] = useState('');
    const [siret, setSiret] = useState('');
    const [category, setCategory] = useState('');
    const [address, setAddress] = useState('');
    const [region, setRegion] = useState('');
    const [contactEmail, setContactEmail] = useState('');
    const [password, setPassword] = useState('');

    const register = useRegisterPartner();

    const handleSubmit = async (e: React.SyntheticEvent<HTMLFormElement>) => {
        e.preventDefault();
        
    }
}
