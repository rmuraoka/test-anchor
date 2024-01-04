import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

// ローカルの翻訳ファイルをインポート
import translationEN from './locales/en/translation.json';
import translationJP from './locales/ja/translation.json';

// 翻訳リソースの設定
const resources = {
    en: {
        translation: translationEN
    },
    ja: {
        translation: translationJP
    }
};

i18n
    // 自動言語検出のためのプラグイン
    .use(LanguageDetector)
    // react-i18next初期化
    .use(initReactI18next)
    .init({
        resources,
        fallbackLng: 'en', // デフォルトの言語
        lng: 'en', // 現在の言語
        keySeparator: false, // キーにドット文字を使用するか
        interpolation: {
            escapeValue: false, // XSS対策を無効化（Reactがデフォルトでエスケープするため）
        },
    });

export default i18n;
