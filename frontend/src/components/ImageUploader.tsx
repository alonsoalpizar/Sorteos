import { useState, useCallback, useRef, useEffect } from 'react';
import { Upload, X, Image as ImageIcon, Star } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { cn } from '@/lib/utils';
import { useUploadImage, useDeleteImage, useSetPrimaryImage } from '@/hooks/useImages';
import { useRaffleDetail } from '@/hooks/useRaffles';

interface ImageUploaderProps {
  raffleId: number;
  maxImages?: number;
  disabled?: boolean;
  className?: string;
}

export function ImageUploader({
  raffleId,
  maxImages = 5,
  disabled = false,
  className,
}: ImageUploaderProps) {
  const [dragActive, setDragActive] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const uploadMutation = useUploadImage();
  const deleteMutation = useDeleteImage();
  const setPrimaryMutation = useSetPrimaryImage();

  // Obtener las imágenes directamente del query para que se actualice automáticamente
  const { data } = useRaffleDetail(raffleId, { includeImages: true });
  const images = data?.images || [];

  // Debug: Log cuando cambien las imágenes
  useEffect(() => {
    console.log('ImageUploader - images changed:', images.length, images);
  }, [images]);

  const handleDrag = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  }, []);

  const validateFile = (file: File): string | null => {
    // Validar tipo
    const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp', 'image/gif'];
    if (!validTypes.includes(file.type)) {
      return 'Tipo de archivo no permitido. Solo se permiten imágenes JPG, PNG, WebP y GIF.';
    }

    // Validar tamaño (10 MB)
    const maxSize = 10 * 1024 * 1024;
    if (file.size > maxSize) {
      return 'El archivo excede el tamaño máximo de 10 MB.';
    }

    return null;
  };

  const handleFiles = (files: FileList | null) => {
    if (!files || files.length === 0) return;

    const file = files[0];
    const error = validateFile(file);

    if (error) {
      alert(error);
      return;
    }

    // Verificar límite de imágenes
    if (images.length >= maxImages) {
      alert(`Máximo ${maxImages} imágenes permitidas por sorteo.`);
      return;
    }

    uploadMutation.mutate({ raffleId, file });
  };

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    if (disabled) return;

    handleFiles(e.dataTransfer.files);
  }, [disabled, raffleId, images.length, maxImages]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    e.preventDefault();
    if (disabled) return;
    handleFiles(e.target.files);
  };

  const handleClick = () => {
    if (disabled) return;
    fileInputRef.current?.click();
  };

  const handleDelete = (imageId: number) => {
    if (confirm('¿Estás seguro de eliminar esta imagen?')) {
      deleteMutation.mutate({ raffleId, imageId });
    }
  };

  const handleSetPrimary = (imageId: number) => {
    setPrimaryMutation.mutate({ raffleId, imageId });
  };

  const sortedImages = [...images].sort((a, b) => {
    if (a.is_primary) return -1;
    if (b.is_primary) return 1;
    return a.display_order - b.display_order;
  });

  return (
    <div className={cn('space-y-4', className)}>
      {/* Upload area */}
      {images.length < maxImages && (
        <div
          className={cn(
            'relative border-2 border-dashed rounded-lg p-8 text-center transition-colors',
            dragActive
              ? 'border-primary bg-primary/5'
              : 'border-gray-300 hover:border-gray-400',
            disabled && 'opacity-50 cursor-not-allowed'
          )}
          onDragEnter={handleDrag}
          onDragLeave={handleDrag}
          onDragOver={handleDrag}
          onDrop={handleDrop}
          onClick={handleClick}
        >
          <input
            ref={fileInputRef}
            type="file"
            className="hidden"
            accept="image/jpeg,image/jpg,image/png,image/webp,image/gif"
            onChange={handleChange}
            disabled={disabled}
          />

          <div className="flex flex-col items-center gap-2">
            <Upload className="w-12 h-12 text-gray-400" />
            <p className="text-sm font-medium text-gray-700">
              Arrastra una imagen o haz clic para seleccionar
            </p>
            <p className="text-xs text-gray-500">
              JPG, PNG, WebP o GIF (máx. 10 MB)
            </p>
            <p className="text-xs text-gray-500">
              {images.length}/{maxImages} imágenes
            </p>
          </div>

          {uploadMutation.isPending && (
            <div className="absolute inset-0 bg-white/80 rounded-lg flex items-center justify-center">
              <div className="flex flex-col items-center gap-2">
                <div className="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin" />
                <p className="text-sm font-medium text-gray-700">Subiendo imagen...</p>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Images grid */}
      {sortedImages.length > 0 && (
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {sortedImages.map((image) => (
            <div
              key={image.id}
              className={cn(
                'relative group aspect-square rounded-lg overflow-hidden border-2 transition-all',
                image.is_primary
                  ? 'border-yellow-400 ring-2 ring-yellow-400 ring-offset-2'
                  : 'border-gray-200 hover:border-gray-300'
              )}
            >
              {/* Image */}
              <img
                src={image.url_medium || image.url_original || ''}
                alt={image.alt_text || 'Imagen del sorteo'}
                className="w-full h-full object-cover"
              />

              {/* Primary badge */}
              {image.is_primary && (
                <div className="absolute top-2 left-2 bg-yellow-400 text-white px-2 py-1 rounded-md text-xs font-medium flex items-center gap-1">
                  <Star className="w-3 h-3 fill-current" />
                  Principal
                </div>
              )}

              {/* Actions overlay */}
              <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
                {!image.is_primary && (
                  <Button
                    size="sm"
                    variant="secondary"
                    onClick={() => handleSetPrimary(image.id)}
                    disabled={setPrimaryMutation.isPending || disabled}
                  >
                    <Star className="w-4 h-4" />
                  </Button>
                )}
                <Button
                  size="sm"
                  variant="destructive"
                  onClick={() => handleDelete(image.id)}
                  disabled={deleteMutation.isPending || disabled}
                >
                  <X className="w-4 h-4" />
                </Button>
              </div>

              {/* Loading overlay */}
              {(deleteMutation.isPending || setPrimaryMutation.isPending) && (
                <div className="absolute inset-0 bg-white/80 flex items-center justify-center">
                  <div className="w-6 h-6 border-4 border-primary border-t-transparent rounded-full animate-spin" />
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Empty state */}
      {sortedImages.length === 0 && (
        <div className="text-center py-8 text-gray-500">
          <ImageIcon className="w-12 h-12 mx-auto mb-2 text-gray-400" />
          <p className="text-sm">No hay imágenes. Sube la primera imagen de tu sorteo.</p>
        </div>
      )}
    </div>
  );
}
