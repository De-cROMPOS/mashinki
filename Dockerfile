FROM libretranslate/libretranslate:latest

# setting up langs
ENV LT_LOAD_ONLY=zh,en,ru

ENV LT_MAX_MEMORY_MB=500

ENV LT_MODEL_SIZE=tiny

EXPOSE 5000 