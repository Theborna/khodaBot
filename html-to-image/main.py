from html2image import Html2Image
hti = Html2Image(
    browser_executable='C:/Users/borna/Downloads/chromiun/chrome-win/chrome-win/chrome.exe',
    size=[300, 300])

hti.screenshot(html_file="temp/latex/ltx.html",
               save_as='latex.png')
