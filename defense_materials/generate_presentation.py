from pathlib import Path

from pptx import Presentation
from pptx.dml.color import RGBColor
from pptx.enum.text import MSO_ANCHOR, PP_ALIGN
from pptx.util import Inches, Pt


OUT = Path(__file__).with_name("defense_presentation.pptx")

NAVY = RGBColor(18, 43, 76)
BLUE = RGBColor(36, 91, 145)
LIGHT_BLUE = RGBColor(232, 241, 250)
GRAY = RGBColor(92, 99, 112)
WHITE = RGBColor(255, 255, 255)
BLACK = RGBColor(24, 28, 35)


def set_run(run, size=20, bold=False, color=BLACK):
    run.font.name = "Arial"
    run.font.size = Pt(size)
    run.font.bold = bold
    run.font.color.rgb = color


def set_fill(shape, color):
    shape.fill.solid()
    shape.fill.fore_color.rgb = color
    shape.line.color.rgb = color


def add_footer(slide, number):
    box = slide.shapes.add_textbox(Inches(0.55), Inches(7.05), Inches(12.2), Inches(0.22))
    p = box.text_frame.paragraphs[0]
    p.text = f"Сарычев В. В. | Магистерская диссертация | {number}/14"
    p.alignment = PP_ALIGN.RIGHT
    set_run(p.runs[0], size=8, color=GRAY)


def add_title(slide, title, subtitle=None, number=None):
    bar = slide.shapes.add_shape(1, Inches(0), Inches(0), Inches(13.333), Inches(0.62))
    set_fill(bar, NAVY)
    title_box = slide.shapes.add_textbox(Inches(0.55), Inches(0.13), Inches(12.2), Inches(0.4))
    p = title_box.text_frame.paragraphs[0]
    p.text = title
    set_run(p.runs[0], size=18, bold=True, color=WHITE)
    if subtitle:
        sub = slide.shapes.add_textbox(Inches(0.75), Inches(0.88), Inches(11.9), Inches(0.45))
        p = sub.text_frame.paragraphs[0]
        p.text = subtitle
        set_run(p.runs[0], size=14, color=GRAY)
    if number:
        add_footer(slide, number)


def add_bullets(slide, items, left=0.85, top=1.35, width=11.7, height=5.55, size=19):
    box = slide.shapes.add_textbox(Inches(left), Inches(top), Inches(width), Inches(height))
    tf = box.text_frame
    tf.clear()
    tf.word_wrap = True
    for i, item in enumerate(items):
        level = 1 if item.startswith("  ") else 0
        text = item.strip()
        p = tf.paragraphs[0] if i == 0 else tf.add_paragraph()
        p.text = text
        p.level = level
        p.space_after = Pt(7 if level == 0 else 3)
        p.line_spacing = 1.05
        p.font.name = "Arial"
        p.font.size = Pt(size - 2 if level else size)
        p.font.color.rgb = BLACK if not text.startswith("Вывод:") else BLUE
        if text.startswith(("Цель:", "Объект:", "Предмет:", "Проблема:", "Вывод:")):
            p.font.bold = True
    return box


def add_highlight(slide, text, left=0.75, top=6.22, width=11.85, height=0.5):
    shape = slide.shapes.add_shape(1, Inches(left), Inches(top), Inches(width), Inches(height))
    set_fill(shape, LIGHT_BLUE)
    shape.text_frame.vertical_anchor = MSO_ANCHOR.MIDDLE
    p = shape.text_frame.paragraphs[0]
    p.text = text
    p.alignment = PP_ALIGN.CENTER
    set_run(p.runs[0], size=15, bold=True, color=NAVY)


def add_table(slide, headers, rows, left=0.5, top=1.35, width=12.35, height=4.75, font_size=11):
    table_shape = slide.shapes.add_table(len(rows) + 1, len(headers), Inches(left), Inches(top), Inches(width), Inches(height))
    table = table_shape.table
    for col_idx, header in enumerate(headers):
        cell = table.cell(0, col_idx)
        cell.text = header
        cell.fill.solid()
        cell.fill.fore_color.rgb = NAVY
        p = cell.text_frame.paragraphs[0]
        p.alignment = PP_ALIGN.CENTER
        p.font.name = "Arial"
        p.font.size = Pt(font_size)
        p.font.bold = True
        p.font.color.rgb = WHITE
    for row_idx, row in enumerate(rows, start=1):
        for col_idx, value in enumerate(row):
            cell = table.cell(row_idx, col_idx)
            cell.text = str(value)
            cell.fill.solid()
            cell.fill.fore_color.rgb = RGBColor(248, 250, 252) if row_idx % 2 else WHITE
            p = cell.text_frame.paragraphs[0]
            p.alignment = PP_ALIGN.CENTER if col_idx else PP_ALIGN.LEFT
            p.font.name = "Arial"
            p.font.size = Pt(font_size)
            p.font.color.rgb = BLACK
    return table_shape


def create_presentation():
    prs = Presentation()
    prs.slide_width = Inches(13.333)
    prs.slide_height = Inches(7.5)
    blank = prs.slide_layouts[6]

    # 1
    slide = prs.slides.add_slide(blank)
    bg = slide.shapes.add_shape(1, Inches(0), Inches(0), Inches(13.333), Inches(7.5))
    set_fill(bg, NAVY)
    title = slide.shapes.add_textbox(Inches(0.85), Inches(1.15), Inches(11.8), Inches(1.45))
    tf = title.text_frame
    tf.word_wrap = True
    p = tf.paragraphs[0]
    p.text = "Исследование и разработка подходов к совершенствованию бизнес-процессов цифровых платформ"
    set_run(p.runs[0], size=31, bold=True, color=WHITE)
    sub = slide.shapes.add_textbox(Inches(0.9), Inches(2.8), Inches(11.5), Inches(0.75))
    p = sub.text_frame.paragraphs[0]
    p.text = "на примере игрового цикла онлайн-казино"
    set_run(p.runs[0], size=22, color=LIGHT_BLUE)
    info = slide.shapes.add_textbox(Inches(0.9), Inches(4.45), Inches(11.7), Inches(1.45))
    for i, line in enumerate([
        "Магистерская диссертация",
        "Сарычев Владислав Витальевич, ПММ-41",
        "НГТУ, факультет прикладной математики и информатики",
        "Руководитель: Стасышин В. М., к.т.н., доцент",
    ]):
        p = info.text_frame.paragraphs[0] if i == 0 else info.text_frame.add_paragraph()
        p.text = line
        set_run(p.runs[0], size=17, color=WHITE if i == 0 else LIGHT_BLUE)

    slides = [
        ("2. Актуальность и проблема", None, [
            "Высокая и неравномерная нагрузка цифровых платформ",
            "Обработка операций в режиме реального времени",
            "Необходимость баланса: производительность, отказы, доходность, anti-fraud",
            "Традиционный анализ недостаточен для стохастических высоконагруженных систем",
            "Проблема: нужен инструмент оценки и оптимизации процесса до внедрения изменений",
        ], None),
        ("3. Цель, объект, предмет и задачи", None, [
            "Цель: разработать подходы к совершенствованию бизнес-процессов цифровых платформ",
            "Объект: бизнес-процессы высоконагруженных цифровых платформ",
            "Предмет: методы моделирования и оптимизации игрового цикла онлайн-казино",
            "Задачи: анализ предметной области, формализация цикла, разработка модели",
            "Реализация программной системы, эксперименты и практические рекомендации",
        ], None),
        ("4. Новизна и практическая значимость", None, [
            "Комплексная модель игрового цикла: нагрузка, очереди, отказы, доходность, anti-fraud",
            "Игровой цикл представлен как система состояний, переходов и событий",
            "Оптимизация параметров основана на результатах имитационного моделирования",
            "Разработана веб-система для проведения экспериментов",
            "Результаты применимы к другим цифровым системам реального времени",
        ], None),
        ("5. Формализация игрового цикла", None, [
            "Состояния: ставка -> проверка параметров -> anti-fraud -> игровая логика -> результат",
            "Дополнительное состояние: ошибка или отказ",
            "Модель: марковский процесс с вероятностными переходами",
            "События: поступление ставки, начало обработки, завершение проверки, результат, ошибка",
            "Формализация позволяет выявлять узкие места процесса",
        ], None),
        ("6. Модель и критерии эффективности", None, [
            "Входные параметры: X = (lambda, mu, k, a, L, F, N, seed)",
            "lambda — интенсивность потока; mu — скорость обработки; k — каналы обработки",
            "a — алгоритмический коэффициент; L — лимит ставки; F — anti-fraud параметр",
            "Критерии: время обработки, очередь, rho, отказы, доход за единицу времени",
            "Anti-fraud метрики: detection rate, false positive rate, false negative rate",
        ], None),
        ("7. Постановка задачи оптимизации", None, [
            "Управляемые параметры: theta = (k, a, L, F)",
            "Цель: максимизировать доход за единицу времени",
            "Ограничения: rho ниже критического уровня и время обработки ниже порога",
            "Сохранение приемлемого качества anti-fraud фильтрации",
            "Методы: Random Search, Grid Search, Genetic Algorithm, Adaptive Optimizer",
        ], None),
        ("8. Программная реализация", None, [
            "Клиент-серверная веб-система моделирования и оптимизации",
            "Backend на Go: симулятор, генератор потока, расчет метрик, anti-fraud, оптимизация, REST API",
            "Frontend на React + TypeScript: параметры, сценарии нагрузки, запуск экспериментов",
            "Визуализация результатов: карточки метрик, таблицы, графики",
            "Экспорт результатов в CSV и JSON; воспроизводимость через seed",
        ], None),
        ("9. Вычислительные эксперименты", None, [
            "1. Влияние интенсивности входящего потока",
            "2. Влияние числа каналов обработки",
            "3. Влияние параметра anti-fraud фильтрации",
            "4. Сравнение baseline и optimized режимов",
            "5. Сравнение методов оптимизации",
        ], None),
    ]

    for idx, (title, subtitle, bullets, highlight) in enumerate(slides, start=2):
        slide = prs.slides.add_slide(blank)
        add_title(slide, title, subtitle, idx)
        add_bullets(slide, bullets)
        if highlight:
            add_highlight(slide, highlight)

    # 10
    slide = prs.slides.add_slide(blank)
    add_title(slide, "10. Эксперимент 1: влияние нагрузки", None, 10)
    add_table(slide, ["Сценарий", "lambda", "rho", "Очередь", "Отказы", "Доход/время"], [
        ["Низкая", "4", "0,2859", "0,0045", "0,0961", "6,7432"],
        ["Средняя", "8", "0,5671", "0,0336", "0,1343", "6,3704"],
        ["Высокая", "12", "0,8528", "0,1926", "0,1739", "-2,3036"],
        ["Стресс", "16", "1,1423", "44,1167", "0,2092", "-14,2652"],
    ], top=1.35, height=3.15, font_size=13)
    add_highlight(slide, "При rho > 1 очередь растет скачкообразно, а доходность становится отрицательной.", top=5.3)

    # 11
    slide = prs.slides.add_slide(blank)
    add_title(slide, "11. Эксперимент 2: масштабирование каналов", "Высокая нагрузка: lambda = 12, mu = 5, a = 1, L = 100, F = 0,25", 11)
    add_table(slide, ["k", "rho", "Очередь", "Отказы", "Доход/время"], [
        ["2", "1,2642", "110,7586", "0,2234", "-12,7840"],
        ["3", "0,8528", "0,1926", "0,1739", "-2,3036"],
        ["4", "0,6476", "0,0316", "0,1487", "5,3741"],
        ["5", "0,5191", "0,0092", "0,1304", "10,4924"],
        ["6", "0,4334", "0,0025", "0,1189", "13,6527"],
    ], top=1.55, height=3.65, font_size=13)
    add_highlight(slide, "Минимально приемлемый режим для сценария высокой нагрузки начинается с k = 4.", top=5.85)

    # 12
    slide = prs.slides.add_slide(blank)
    add_title(slide, "12. Эксперимент 3: anti-fraud параметр", None, 12)
    add_table(slide, ["F", "Detection", "FNR", "Время", "Отказы", "Доход/время"], [
        ["0,0", "0,8996", "0,1004", "0,2011", "0,1268", "7,6495"],
        ["0,1", "0,9844", "0,0156", "0,2068", "0,1320", "6,7241"],
        ["0,2", "0,9974", "0,0026", "0,2127", "0,1343", "6,3704"],
        ["0,3", "1,0000", "0,0000", "0,2187", "0,1398", "5,2159"],
        ["0,5", "1,0000", "0,0000", "0,2303", "0,1524", "2,8466"],
    ], top=1.35, height=3.65, font_size=12)
    add_highlight(slide, "Рациональная область F = 0,2-0,3: безопасность растет без чрезмерной потери эффективности.", top=5.78)

    # 13
    slide = prs.slides.add_slide(blank)
    add_title(slide, "13. Эксперимент 4: baseline vs optimized", "Исходный сценарий высокой нагрузки: lambda = 12, mu = 5, k = 3, a = 1, L = 100, F = 0,25", 13)
    add_table(slide, ["Метрика", "Baseline", "Optimized"], [
        ["k", "3", "5"],
        ["a", "1", "1,5"],
        ["L", "100", "150"],
        ["F", "0,25", "0,15"],
        ["Время обработки", "0,2132", "0,1410"],
        ["Время ожидания", "0,1926", "0,0009"],
        ["rho", "0,8528", "0,3383"],
        ["Вероятность отказа", "0,1739", "0,0956"],
        ["Доход/время", "-2,3036", "42,0042"],
        ["Detection rate", "0,9987", "0,9988"],
    ], top=1.45, height=4.75, font_size=10)
    add_highlight(slide, "Оптимизация улучшила производительность и доходность без ухудшения anti-fraud качества.", top=6.47, height=0.42)

    # 14
    slide = prs.slides.add_slide(blank)
    add_title(slide, "14. Сравнение методов и итоговые выводы", None, 14)
    add_table(slide, ["Метод", "Итерации", "Доход/время", "Очередь", "rho", "Fraud Detection"], [
        ["Random Search", "100", "40,0864", "0,0012", "0,3478", "1,0000"],
        ["Grid Search", "320", "42,0042", "0,0009", "0,3383", "0,9988"],
        ["Genetic Algorithm", "100", "42,0042", "0,0009", "0,3383", "0,9988"],
        ["Adaptive Optimizer", "100", "36,6021", "0,0033", "0,4224", "0,9962"],
    ], top=1.15, height=2.75, font_size=11)
    add_bullets(slide, [
        "Цель работы достигнута: разработана модель и программная система",
        "Grid Search и Genetic Algorithm показали лучшие результаты",
        "Genetic Algorithm достиг результата полного перебора при меньшем числе проверок",
        "Практические рекомендации: контролировать rho, масштабировать k, настраивать F как компромисс",
    ], top=4.25, height=1.8, size=16)
    add_highlight(slide, "Имитационное моделирование и оптимизация применимы для совершенствования бизнес-процессов цифровых платформ.", top=6.35, height=0.48)

    prs.save(OUT)


if __name__ == "__main__":
    create_presentation()
    print(f"Saved: {OUT}")
