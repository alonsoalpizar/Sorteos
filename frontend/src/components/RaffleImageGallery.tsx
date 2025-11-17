import { useState } from 'react';
import { X } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { RaffleImage } from '@/types/raffle';

interface RaffleImageGalleryProps {
  images: RaffleImage[];
  className?: string;
}

export function RaffleImageGallery({ images, className }: RaffleImageGalleryProps) {
  const [selectedIndex, setSelectedIndex] = useState(0);
  const [lightboxOpen, setLightboxOpen] = useState(false);

  // Ordenar imágenes: primaria primero, luego por display_order
  const sortedImages = [...images].sort((a, b) => {
    if (a.is_primary) return -1;
    if (b.is_primary) return 1;
    return a.display_order - b.display_order;
  });

  if (sortedImages.length === 0) {
    return null;
  }

  const currentImage = sortedImages[selectedIndex];

  return (
    <div className={cn('space-y-4', className)}>
      {/* Imagen principal */}
      <div
        className="relative max-w-xs mx-auto bg-gray-100 dark:bg-gray-800 rounded-lg overflow-hidden cursor-pointer group"
        onClick={() => setLightboxOpen(true)}
      >
        <img
          src={currentImage.url_large || currentImage.url_original || ''}
          alt={currentImage.alt_text || 'Imagen del sorteo'}
          className="w-full h-auto object-contain transition-transform duration-300 group-hover:scale-105"
        />

        {/* Overlay de zoom en hover */}
        <div className="absolute inset-0 bg-black/0 group-hover:bg-black/20 transition-all flex items-center justify-center">
          <div className="opacity-0 group-hover:opacity-100 transition-opacity bg-white/90 dark:bg-gray-800/90 px-4 py-2 rounded-full text-sm font-medium">
            Click para ampliar
          </div>
        </div>

        {/* Badge de imagen principal */}
        {currentImage.is_primary && (
          <div className="absolute top-2 right-2 bg-blue-600 text-white px-2 py-0.5 rounded-full text-xs font-medium">
            Principal
          </div>
        )}

        {/* Contador de imágenes */}
        {sortedImages.length > 1 && (
          <div className="absolute bottom-2 right-2 bg-black/60 text-white px-2 py-0.5 rounded-full text-xs font-medium">
            {selectedIndex + 1} / {sortedImages.length}
          </div>
        )}
      </div>

      {/* Thumbnails (solo si hay más de 1 imagen) */}
      {sortedImages.length > 1 && (
        <div className="flex gap-1.5 overflow-x-auto pb-2 justify-center">
          {sortedImages.map((image, index) => (
            <button
              key={image.id}
              onClick={() => setSelectedIndex(index)}
              className={cn(
                'relative flex-shrink-0 w-12 h-12 rounded-md overflow-hidden border-2 transition-all',
                index === selectedIndex
                  ? 'border-blue-600 ring-2 ring-blue-600 ring-offset-1'
                  : 'border-gray-200 dark:border-gray-700 hover:border-gray-400 dark:hover:border-gray-500'
              )}
            >
              <img
                src={image.url_thumbnail || image.url_medium || ''}
                alt={image.alt_text || `Imagen ${index + 1}`}
                className="w-full h-full object-cover"
              />

              {/* Overlay en thumbnails no seleccionados */}
              {index !== selectedIndex && (
                <div className="absolute inset-0 bg-black/30" />
              )}
            </button>
          ))}
        </div>
      )}

      {/* Lightbox */}
      {lightboxOpen && (
        <div
          className="fixed inset-0 z-50 bg-black/95 flex items-center justify-center p-4"
          onClick={() => setLightboxOpen(false)}
        >
          {/* Botón cerrar */}
          <button
            className="absolute top-4 right-4 p-2 bg-white/10 hover:bg-white/20 rounded-full transition-colors"
            onClick={() => setLightboxOpen(false)}
          >
            <X className="w-6 h-6 text-white" />
          </button>

          {/* Contador */}
          {sortedImages.length > 1 && (
            <div className="absolute top-4 left-1/2 -translate-x-1/2 bg-white/10 text-white px-4 py-2 rounded-full text-sm font-medium">
              {selectedIndex + 1} / {sortedImages.length}
            </div>
          )}

          {/* Imagen en tamaño completo */}
          <div className="relative max-w-7xl max-h-full">
            <img
              src={currentImage.url_original || currentImage.url_large || ''}
              alt={currentImage.alt_text || 'Imagen del sorteo'}
              className="max-w-full max-h-[90vh] object-contain"
              onClick={(e) => e.stopPropagation()}
            />
          </div>

          {/* Navegación (flechas) si hay múltiples imágenes */}
          {sortedImages.length > 1 && (
            <>
              {/* Flecha izquierda */}
              <button
                className="absolute left-4 top-1/2 -translate-y-1/2 p-3 bg-white/10 hover:bg-white/20 rounded-full transition-colors disabled:opacity-50"
                onClick={(e) => {
                  e.stopPropagation();
                  setSelectedIndex((prev) => (prev > 0 ? prev - 1 : sortedImages.length - 1));
                }}
                disabled={sortedImages.length <= 1}
              >
                <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                </svg>
              </button>

              {/* Flecha derecha */}
              <button
                className="absolute right-4 top-1/2 -translate-y-1/2 p-3 bg-white/10 hover:bg-white/20 rounded-full transition-colors disabled:opacity-50"
                onClick={(e) => {
                  e.stopPropagation();
                  setSelectedIndex((prev) => (prev < sortedImages.length - 1 ? prev + 1 : 0));
                }}
                disabled={sortedImages.length <= 1}
              >
                <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </button>
            </>
          )}

          {/* Thumbnails en lightbox */}
          {sortedImages.length > 1 && (
            <div className="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2 max-w-full overflow-x-auto px-4">
              {sortedImages.map((image, index) => (
                <button
                  key={image.id}
                  onClick={(e) => {
                    e.stopPropagation();
                    setSelectedIndex(index);
                  }}
                  className={cn(
                    'flex-shrink-0 w-16 h-16 rounded-lg overflow-hidden border-2 transition-all',
                    index === selectedIndex
                      ? 'border-white ring-2 ring-white'
                      : 'border-white/30 hover:border-white/60'
                  )}
                >
                  <img
                    src={image.url_thumbnail || image.url_medium || ''}
                    alt={`Imagen ${index + 1}`}
                    className="w-full h-full object-cover"
                  />
                </button>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
