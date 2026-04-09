const headerOffset = () => {
    const header = document.querySelector('header');
    return header ? header.getBoundingClientRect().height + 16 : 96;
};

const scrollToTarget = (target) => {
    if (!target) return;

    if (target.id === 'inicio') {
        window.scrollTo({ top: 0, left: 0, behavior: 'smooth' });
        return;
    }

    const anchor = target.matches('section')
        ? (target.querySelector('[data-scroll-anchor]') || target.querySelector('h2, h3, h4') || target)
        : target;

    const top = window.scrollY + anchor.getBoundingClientRect().top - headerOffset();
    window.scrollTo({ top: Math.max(0, top), left: 0, behavior: 'smooth' });
};

document.addEventListener('click', (event) => {
    const link = event.target.closest('a[href^="#"]');
    if (!link) return;

    const href = link.getAttribute('href');
    if (!href || href === '#') return;

    const target = document.querySelector(href);
    if (!target) return;

    event.preventDefault();
    scrollToTarget(target);
});

// Função para alternar todas as galerias
function toggleAllGalleries(button) {
    const allGalleries = document.querySelectorAll('[data-gallery]');
    const toggleText = button.querySelector('.toggle-text');
    const chevron = button.querySelector('.fa-chevron-down');
    const animationMs = 220;

    const centerButtonInViewport = () => {
        const buttonRect = button.getBoundingClientRect();
        const targetY = window.scrollY + buttonRect.top - (window.innerHeight / 2) + (buttonRect.height / 2);
        const maxScroll = Math.max(0, document.documentElement.scrollHeight - window.innerHeight);
        const clampedY = Math.min(Math.max(0, targetY), maxScroll);

        window.scrollTo({ top: clampedY, left: 0, behavior: 'smooth' });
    };

    const keepButtonStable = (initialTop, durationMs, onDone) => {
        const startedAt = performance.now();

        const tick = (now) => {
            const delta = button.getBoundingClientRect().top - initialTop;
            if (delta !== 0) {
                window.scrollTo({ top: Math.max(0, window.scrollY + delta), left: 0, behavior: 'auto' });
            }

            if (now - startedAt < durationMs) {
                requestAnimationFrame(tick);
                return;
            }

            onDone();
        };

        requestAnimationFrame(tick);
    };

    if (button.dataset.busy === 'true') return;
    button.dataset.busy = 'true';

    // Verificar estado atual do botão
    const isExpanded = button.dataset.expanded === 'true';

    if (isExpanded) {
        const buttonTopBefore = button.getBoundingClientRect().top;

        // Esconder imagens - voltar ao estado inicial
        allGalleries.forEach((gallery) => {
            const extraItems = gallery.querySelectorAll('.gallery-item-extra');
            extraItems.forEach((item) => {
                item.classList.remove('is-visible');
            });
        });

        toggleText.textContent = 'Ver Mais Imagens';
        chevron.style.transform = 'rotate(0deg)';
        button.dataset.expanded = 'false';

        keepButtonStable(buttonTopBefore, animationMs, () => {
            allGalleries.forEach((gallery) => {
                const extraItems = gallery.querySelectorAll('.gallery-item-extra');
                extraItems.forEach((item) => {
                    item.style.display = 'none';
                });
            });

            // Garante estabilidade final apos remover os itens do fluxo.
            const buttonTopAfter = button.getBoundingClientRect().top;
            const delta = buttonTopAfter - buttonTopBefore;
            if (delta !== 0) {
                const targetY = Math.max(0, window.scrollY + delta);
                window.scrollTo({ top: targetY, left: 0, behavior: 'auto' });
            }

            requestAnimationFrame(() => {
                centerButtonInViewport();
            });

            button.dataset.busy = 'false';
        });
    } else {
        // Mostrar imagens
        allGalleries.forEach((gallery) => {
            const extraItems = gallery.querySelectorAll('.gallery-item-extra');
            extraItems.forEach((item) => {
                item.style.display = '';
            });

            requestAnimationFrame(() => {
                extraItems.forEach((item) => {
                    item.classList.add('is-visible');
                });
            });
        });

        toggleText.textContent = 'Ver Menos Imagens';
        chevron.style.transform = 'rotate(180deg)';
        button.dataset.expanded = 'true';
        setTimeout(() => {
            button.dataset.busy = 'false';
        }, animationMs);
    }
}

const getIconClass = (name) => {
    const n = name.toLowerCase();
    if (n.includes('mecânica pesada')) return 'fa-truck';
    if (n.includes('diferencial')) return 'fa-gears';
    if (n.includes('motor')) return 'fa-cogs';
    if (n.includes('câmbio')) return 'fa-gear';
    if (n.includes('elétrica')) return 'fa-bolt';
    if (n.includes('módulo') || n.includes('automação')) return 'fa-microchip';
    if (n.includes('arla') || n.includes('injeção') || n.includes('common rail')) return 'fa-wind';
    if (n.includes('rastreamento') || n.includes('diagnóstico') || n.includes('raster')) return 'fa-laptop-code';
    if (n.includes('borracharia')) return 'fa-dharmachakra';
    if (n.includes('solda')) return 'fa-fire';
    if (n.includes('torno')) return 'fa-screwdriver-wrench';
    return 'fa-wrench';
};

fetch('/api/services')
    .then((res) => res.json())
    .then((services) => {
        const container = document.getElementById('services-container');
        container.innerHTML = services.map((service) => `
            <div class="service-card bg-white p-8 border border-gray-100 shadow-sm hover:shadow-2xl transition-all group relative overflow-hidden">
                <div class="absolute top-0 right-0 w-16 h-16 bg-shop-red/5 -mr-8 -mt-8 rounded-full group-hover:bg-shop-red/10 transition-all"></div>
                <div class="text-shop-red text-3xl mb-6 relative z-10">
                    <i class="fas ${getIconClass(service.name)}"></i>
                </div>
                <h4 class="text-lg font-black text-shop-gray-dark uppercase italic tracking-tighter mb-3 relative z-10">${service.name}</h4>
                <p class="text-sm text-gray-500 font-medium leading-relaxed relative z-10">${service.description}</p>
                <div class="mt-6 border-t border-gray-50 pt-6 opacity-0 group-hover:opacity-100 transition-all">
                    <a href="#contato" class="btn-clean text-shop-red text-[10px] font-black uppercase tracking-widest items-center gap-2 hover:opacity-80">
                        Saiba Mais <i class="fas fa-chevron-right"></i>
                    </a>
                </div>
            </div>
        `).join('');
    })
    .catch((err) => console.error('Erro:', err));

const logoWrapper = document.getElementById('columbina-wrapper');
const floatingLogo = document.getElementById('imagemFlutuante');

if (logoWrapper && floatingLogo) {
    let animationFrame = null;
    let targetX = 0;
    const baseY = -4.0742;
    let targetY = baseY;
    let currentX = 0;
    let currentY = baseY;
    let floatTime = 0;

    const maxOffset = 8;
    const followSpeed = 0.21;
    const floatAmplitude = 2.2;
    const floatVelocity = 0.003;

    const animateLogo = (time) => {
        floatTime = time;
        currentX += (targetX - currentX) * followSpeed;
        currentY += (targetY - currentY) * followSpeed;
        const slowFloat = Math.sin(floatTime * floatVelocity) * floatAmplitude;

        floatingLogo.style.transform = `translate3d(${currentX}px, ${currentY + slowFloat}px, 0)`;
        animationFrame = requestAnimationFrame(animateLogo);
    };

    const moveLogo = (event) => {
        const rect = logoWrapper.getBoundingClientRect();
        const normalizedX = (event.clientX - rect.left) / rect.width - 0.5;
        const normalizedY = (event.clientY - rect.top) / rect.height - 0.5;

        targetX = normalizedX * maxOffset;
        targetY = baseY + normalizedY * maxOffset;

        if (!animationFrame) {
            animationFrame = requestAnimationFrame(animateLogo);
        }
    };

    const resetLogo = () => {
        targetX = 0;
        targetY = baseY;

        if (!animationFrame) {
            animationFrame = requestAnimationFrame(animateLogo);
        }
    };

    logoWrapper.addEventListener('mousemove', moveLogo);
    logoWrapper.addEventListener('mouseenter', moveLogo);
    logoWrapper.addEventListener('mouseleave', resetLogo);
    animationFrame = requestAnimationFrame(animateLogo);
}
