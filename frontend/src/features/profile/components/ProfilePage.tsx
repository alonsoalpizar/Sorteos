import { useState } from 'react';
import { User, Wallet, FileText, Edit2, Upload, Check, X } from 'lucide-react';
import { Card } from '../../../components/ui/Card';
import { Button } from '../../../components/ui/Button';
import { LoadingSpinner } from '../../../components/ui/LoadingSpinner';
import { useProfile } from '../hooks/useProfile';
import { useUpdateProfile } from '../hooks/useUpdateProfile';
import { useConfigureIBAN } from '../hooks/useConfigureIBAN';
import { useUploadKYCDocument } from '../hooks/useUploadKYCDocument';
import { formatCRC } from '../../../types/wallet';
import { toast } from 'sonner';

export const ProfilePage = () => {
  const { data, isLoading, error } = useProfile();
  const updateProfile = useUpdateProfile();
  const configureIBAN = useConfigureIBAN();
  const uploadKYCDocument = useUploadKYCDocument();

  const [isEditingPersonal, setIsEditingPersonal] = useState(false);
  const [isEditingIBAN, setIsEditingIBAN] = useState(false);

  // Form states
  const [personalForm, setPersonalForm] = useState({
    first_name: '',
    last_name: '',
    phone: '',
    cedula: '',
    date_of_birth: '',
    address_line1: '',
    address_line2: '',
    city: '',
    state: '',
    postal_code: '',
  });

  const [ibanForm, setIbanForm] = useState('');
  const [uploadingDoc, setUploadingDoc] = useState<string | null>(null);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-6">
        <Card className="p-6 bg-red-50 border-red-200">
          <p className="text-red-800">Error al cargar el perfil: {error.message}</p>
        </Card>
      </div>
    );
  }

  if (!data) {
    return null;
  }

  const { user, wallet, kyc_documents, can_withdraw } = data;

  // KYC status
  const kycStatus = {
    cedula_front: kyc_documents.find((d) => d.document_type === 'cedula_front'),
    cedula_back: kyc_documents.find((d) => d.document_type === 'cedula_back'),
    selfie: kyc_documents.find((d) => d.document_type === 'selfie'),
  };

  const kycLevelLabels = {
    none: 'Sin verificar',
    email_verified: 'Email verificado',
    phone_verified: 'Teléfono verificado',
    cedula_verified: 'Cédula verificada',
    full_kyc: 'Verificación completa',
  };

  // Handlers
  const handleEditPersonal = () => {
    setPersonalForm({
      first_name: user.first_name || '',
      last_name: user.last_name || '',
      phone: user.phone || '',
      cedula: user.cedula || '',
      date_of_birth: user.date_of_birth || '',
      address_line1: user.address_line1 || '',
      address_line2: user.address_line2 || '',
      city: user.city || '',
      state: user.state || '',
      postal_code: user.postal_code || '',
    });
    setIsEditingPersonal(true);
  };

  const handleSavePersonal = async () => {
    try {
      await updateProfile.mutateAsync(personalForm);
      setIsEditingPersonal(false);
      toast.success('Perfil actualizado correctamente');
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Error al actualizar perfil');
    }
  };

  const handleCancelPersonal = () => {
    setIsEditingPersonal(false);
  };

  const handleEditIBAN = () => {
    setIbanForm(user.iban || '');
    setIsEditingIBAN(true);
  };

  const handleSaveIBAN = async () => {
    if (!ibanForm || ibanForm.length !== 24) {
      toast.error('IBAN debe tener 24 caracteres (CR + 22 dígitos)');
      return;
    }
    try {
      await configureIBAN.mutateAsync(ibanForm);
      setIsEditingIBAN(false);
      toast.success('IBAN configurado correctamente');
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Error al configurar IBAN');
    }
  };

  const handleCancelIBAN = () => {
    setIsEditingIBAN(false);
  };

  const handleFileUpload = async (
    docType: 'cedula_front' | 'cedula_back' | 'selfie',
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (!file) return;

    // Validar tipo de archivo
    if (!file.type.startsWith('image/')) {
      toast.error('Solo se permiten imágenes');
      return;
    }

    // Validar tamaño (5MB max)
    if (file.size > 5 * 1024 * 1024) {
      toast.error('La imagen no debe superar 5MB');
      return;
    }

    setUploadingDoc(docType);

    try {
      // TODO: Implementar upload real a servidor
      // Por ahora simulamos con un placeholder
      const fakeUrl = `https://sorteos.club/uploads/kyc/${Date.now()}_${file.name}`;

      await uploadKYCDocument.mutateAsync({
        documentType: docType,
        fileUrl: fakeUrl,
      });

      toast.success('Documento subido correctamente');
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Error al subir documento');
    } finally {
      setUploadingDoc(null);
    }
  };

  return (
    <div className="container mx-auto p-6 space-y-6">
      <h1 className="text-3xl font-bold text-slate-900">Mi Perfil</h1>

      {/* Información Personal */}
      <Card className="p-6">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-3">
            <User className="w-6 h-6 text-blue-600" />
            <h2 className="text-xl font-semibold text-slate-900">Información Personal</h2>
          </div>
          {!isEditingPersonal && (
            <Button
              variant="outline"
              size="sm"
              onClick={handleEditPersonal}
              className="flex items-center gap-2"
            >
              <Edit2 className="w-4 h-4" />
              Editar
            </Button>
          )}
        </div>

        {isEditingPersonal ? (
          <div className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">
                  Nombre
                </label>
                <input
                  type="text"
                  value={personalForm.first_name}
                  onChange={(e) => setPersonalForm({ ...personalForm, first_name: e.target.value })}
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">
                  Apellidos
                </label>
                <input
                  type="text"
                  value={personalForm.last_name}
                  onChange={(e) => setPersonalForm({ ...personalForm, last_name: e.target.value })}
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">
                  Teléfono
                </label>
                <input
                  type="tel"
                  value={personalForm.phone}
                  onChange={(e) => setPersonalForm({ ...personalForm, phone: e.target.value })}
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="88887777"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">
                  Cédula
                </label>
                <input
                  type="text"
                  value={personalForm.cedula}
                  onChange={(e) => setPersonalForm({ ...personalForm, cedula: e.target.value })}
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="1-2345-6789"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">
                  Fecha de Nacimiento
                </label>
                <input
                  type="date"
                  value={personalForm.date_of_birth}
                  onChange={(e) => setPersonalForm({ ...personalForm, date_of_birth: e.target.value })}
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Sección de Dirección */}
            <div className="mt-6 pt-6 border-t border-slate-200">
              <h3 className="text-lg font-semibold text-slate-900 mb-4">Dirección</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="md:col-span-2">
                  <label className="block text-sm font-medium text-slate-700 mb-1">
                    Dirección Línea 1
                  </label>
                  <input
                    type="text"
                    value={personalForm.address_line1}
                    onChange={(e) => setPersonalForm({ ...personalForm, address_line1: e.target.value })}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="Calle principal, número de casa"
                  />
                </div>

                <div className="md:col-span-2">
                  <label className="block text-sm font-medium text-slate-700 mb-1">
                    Dirección Línea 2
                  </label>
                  <input
                    type="text"
                    value={personalForm.address_line2}
                    onChange={(e) => setPersonalForm({ ...personalForm, address_line2: e.target.value })}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="Apartamento, suite (opcional)"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">
                    Ciudad
                  </label>
                  <input
                    type="text"
                    value={personalForm.city}
                    onChange={(e) => setPersonalForm({ ...personalForm, city: e.target.value })}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="San José"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">
                    Provincia
                  </label>
                  <input
                    type="text"
                    value={personalForm.state}
                    onChange={(e) => setPersonalForm({ ...personalForm, state: e.target.value })}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="San José"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">
                    Código Postal
                  </label>
                  <input
                    type="text"
                    value={personalForm.postal_code}
                    onChange={(e) => setPersonalForm({ ...personalForm, postal_code: e.target.value })}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="10101"
                  />
                </div>
              </div>
            </div>

            <div className="flex gap-2 justify-end">
              <Button
                variant="outline"
                onClick={handleCancelPersonal}
                className="flex items-center gap-2"
              >
                <X className="w-4 h-4" />
                Cancelar
              </Button>
              <Button
                onClick={handleSavePersonal}
                disabled={updateProfile.isPending}
                className="flex items-center gap-2"
              >
                <Check className="w-4 h-4" />
                Guardar
              </Button>
            </div>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-slate-600">Nombre Completo</p>
              <p className="font-medium text-slate-900">
                {user.first_name && user.last_name
                  ? `${user.first_name} ${user.last_name}`
                  : 'No configurado'}
              </p>
            </div>

            <div>
              <p className="text-sm text-slate-600">Email</p>
              <p className="font-medium text-slate-900">{user.email}</p>
              {user.email_verified && (
                <span className="text-xs text-green-600">✓ Verificado</span>
              )}
            </div>

            <div>
              <p className="text-sm text-slate-600">Teléfono</p>
              <p className="font-medium text-slate-900">{user.phone || 'No configurado'}</p>
              {user.phone_verified && user.phone && (
                <span className="text-xs text-green-600">✓ Verificado</span>
              )}
            </div>

            <div>
              <p className="text-sm text-slate-600">Cédula</p>
              <p className="font-medium text-slate-900">{user.cedula || 'No configurada'}</p>
            </div>

            <div>
              <p className="text-sm text-slate-600">Fecha de Nacimiento</p>
              <p className="font-medium text-slate-900">
                {user.date_of_birth
                  ? new Date(user.date_of_birth).toLocaleDateString('es-CR')
                  : 'No configurada'}
              </p>
            </div>

            <div>
              <p className="text-sm text-slate-600">Nivel KYC</p>
              <p className="font-medium text-slate-900">
                {kycLevelLabels[user.kyc_level]}
              </p>
            </div>
          </div>
        )}
      </Card>

      {/* Estado de Wallet */}
      <Card className="p-6">
        <div className="flex items-center gap-3 mb-4">
          <Wallet className="w-6 h-6 text-green-600" />
          <h2 className="text-xl font-semibold text-slate-900">Billetera</h2>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <p className="text-sm text-slate-600">Saldo Disponible</p>
            <p className="text-2xl font-bold text-blue-600">
              {formatCRC(parseFloat(wallet.balance_available))}
            </p>
            <p className="text-xs text-slate-500">Para comprar tickets</p>
          </div>

          <div>
            <p className="text-sm text-slate-600">Ganancias</p>
            <p className="text-2xl font-bold text-green-600">
              {formatCRC(parseFloat(wallet.earnings_balance))}
            </p>
            <p className="text-xs text-slate-500">De tus sorteos</p>
          </div>

          <div>
            <p className="text-sm text-slate-600">Estado de Retiros</p>
            <p className={`text-lg font-semibold ${can_withdraw ? 'text-green-600' : 'text-orange-600'}`}>
              {can_withdraw ? '✓ Habilitado' : '✗ Deshabilitado'}
            </p>
            <p className="text-xs text-slate-500">
              {can_withdraw
                ? 'Puedes retirar ganancias'
                : 'Completa KYC y configura IBAN'}
            </p>
          </div>
        </div>
      </Card>

      {/* Documentos KYC */}
      <Card className="p-6">
        <div className="flex items-center gap-3 mb-4">
          <FileText className="w-6 h-6 text-amber-600" />
          <h2 className="text-xl font-semibold text-slate-900">Documentos KYC</h2>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {(['cedula_front', 'cedula_back', 'selfie'] as const).map((docType) => {
            const doc = kycStatus[docType];
            const labels = {
              cedula_front: 'Cédula (Frente)',
              cedula_back: 'Cédula (Dorso)',
              selfie: 'Selfie con Cédula',
            };

            const isUploading = uploadingDoc === docType;

            return (
              <div
                key={docType}
                className="p-4 border rounded-lg"
              >
                <p className="font-medium text-slate-900 mb-2">
                  {labels[docType]}
                </p>

                {doc ? (
                  <div className="space-y-2">
                    <span
                      className={`inline-block px-2 py-1 rounded text-xs font-medium ${
                        doc.verification_status === 'approved'
                          ? 'bg-green-100 text-green-800'
                          : doc.verification_status === 'rejected'
                          ? 'bg-red-100 text-red-800'
                          : 'bg-amber-100 text-amber-800'
                      }`}
                    >
                      {doc.verification_status === 'approved'
                        ? 'Aprobado'
                        : doc.verification_status === 'rejected'
                        ? 'Rechazado'
                        : 'Pendiente'}
                    </span>
                    {doc.rejected_reason && (
                      <p className="text-xs text-red-600 mt-1">{doc.rejected_reason}</p>
                    )}
                    {doc.verification_status !== 'approved' && (
                      <label className="block">
                        <input
                          type="file"
                          accept="image/*"
                          onChange={(e) => handleFileUpload(docType, e)}
                          disabled={isUploading}
                          className="hidden"
                        />
                        <Button
                          variant="outline"
                          size="sm"
                          className="w-full flex items-center gap-2"
                          disabled={isUploading}
                          onClick={(e) => {
                            e.preventDefault();
                            (e.currentTarget.previousElementSibling as HTMLInputElement)?.click();
                          }}
                        >
                          <Upload className="w-4 h-4" />
                          {isUploading ? 'Subiendo...' : 'Resubir'}
                        </Button>
                      </label>
                    )}
                  </div>
                ) : (
                  <label className="block">
                    <input
                      type="file"
                      accept="image/*"
                      onChange={(e) => handleFileUpload(docType, e)}
                      disabled={isUploading}
                      className="hidden"
                    />
                    <Button
                      variant="outline"
                      size="sm"
                      className="w-full flex items-center gap-2"
                      disabled={isUploading}
                      onClick={(e) => {
                        e.preventDefault();
                        (e.currentTarget.previousElementSibling as HTMLInputElement)?.click();
                      }}
                    >
                      <Upload className="w-4 h-4" />
                      {isUploading ? 'Subiendo...' : 'Subir'}
                    </Button>
                  </label>
                )}
              </div>
            );
          })}
        </div>

        <div className="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-lg">
          <p className="text-sm text-blue-800">
            <strong>Importante:</strong> Sube imágenes claras y legibles. Los documentos serán revisados por nuestro equipo en un plazo de 24-48 horas.
          </p>
        </div>
      </Card>

      {/* IBAN */}
      <Card className="p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-slate-900">Información Bancaria</h2>
          {!isEditingIBAN && !user.iban && user.kyc_level === 'full_kyc' && (
            <Button
              variant="outline"
              size="sm"
              onClick={handleEditIBAN}
              className="flex items-center gap-2"
            >
              <Edit2 className="w-4 h-4" />
              Configurar
            </Button>
          )}
          {!isEditingIBAN && user.iban && (
            <Button
              variant="outline"
              size="sm"
              onClick={handleEditIBAN}
              className="flex items-center gap-2"
            >
              <Edit2 className="w-4 h-4" />
              Editar
            </Button>
          )}
        </div>

        {user.kyc_level !== 'full_kyc' && !user.iban ? (
          <div className="p-4 bg-amber-50 border border-amber-200 rounded-lg">
            <p className="text-sm text-amber-800">
              Debes completar la verificación KYC antes de configurar tu IBAN para retiros.
            </p>
          </div>
        ) : isEditingIBAN ? (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">
                IBAN de Costa Rica
              </label>
              <input
                type="text"
                value={ibanForm}
                onChange={(e) => setIbanForm(e.target.value.toUpperCase())}
                className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono"
                placeholder="CR12345678901234567890"
                maxLength={24}
              />
              <p className="text-xs text-slate-500 mt-1">
                Formato: CR + 22 dígitos (24 caracteres en total)
              </p>
            </div>

            <div className="flex gap-2 justify-end">
              <Button
                variant="outline"
                onClick={handleCancelIBAN}
                className="flex items-center gap-2"
              >
                <X className="w-4 h-4" />
                Cancelar
              </Button>
              <Button
                onClick={handleSaveIBAN}
                disabled={configureIBAN.isPending}
                className="flex items-center gap-2"
              >
                <Check className="w-4 h-4" />
                Guardar
              </Button>
            </div>
          </div>
        ) : user.iban ? (
          <div>
            <p className="text-sm text-slate-600">IBAN Configurado</p>
            <p className="font-mono text-lg font-medium text-slate-900">{user.iban}</p>
          </div>
        ) : null}
      </Card>
    </div>
  );
};
