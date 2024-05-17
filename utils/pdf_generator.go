package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func GenerateCertificatePdf(certificateData entities.CertificateData) (string, error) {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(20, 10, 20)

	buildHeading(m)
	buildCertificate(m, certificateData)

	fileName := fmt.Sprintf("pdfs/certificate-%d.pdf", certificateData.ConsultationId)
	err := m.OutputFileAndClose(fileName)

	if err != nil {
		return "", err
	}

	return fileName, nil
}

func GeneratePrescriptionPdf(prescriptionData entities.PrescriptionData) (string, error) {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(20, 10, 20)

	buildHeading(m)
	buildPrescription(m, prescriptionData)

	fileName := fmt.Sprintf("pdfs/prescription-%d.pdf", prescriptionData.ConsultationId)
	err := m.OutputFileAndClose(fileName)

	if err != nil {
		return "", err
	}

	return fileName, nil
}

func buildHeading(m pdf.Maroto) {
	m.RegisterHeader(func() {
		m.Row(20, func() {
			m.Col(12, func() {
				err := m.FileImage("images/logo.png", props.Rect{
					Center:  true,
					Percent: 50,
				})

				if err != nil {
					log.Println("Image file was not loaded 😱 - ", err)
				}

			})
		})
	})
}

func buildCertificate(m pdf.Maroto, certificateData entities.CertificateData) {
	headings := getHeadings()
	contents := [][]string{{"Name", fmt.Sprintf(" : %s", certificateData.PatientName)}, {"Gender", fmt.Sprintf(" : %s", certificateData.PatientGender.Name)}, {"Birth Date", fmt.Sprintf(" : %s", formatDate(certificateData.PatientBirthDate))}, {"Age", fmt.Sprintf(" : %d y.o.", certificateData.PatientAge)}}

	m.SetBackgroundColor(getTealColor())
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Patient Data", props.Text{
				Top:    2,
				Size:   13,
				Color:  color.NewWhite(),
				Family: consts.Arial,
				Style:  consts.Bold,
				Align:  consts.Center,
			})
		})
	})

	m.SetBackgroundColor(color.NewWhite())

	m.TableList(headings, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      10,
			GridSizes: []uint{3, 3},
		},
		ContentProp: props.TableListContent{
			Size:      10,
			GridSizes: []uint{3, 3},
		},
		Align:                  consts.Left,
		HeaderContentSpace:     1,
		Line:                   false,
		VerticalContentPadding: 2,
	})

	m.Row(10, func() {
		m.Col(12, func() {
		})
	})

	m.SetBackgroundColor(getTealColor())
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Doctor Note", props.Text{
				Top:    2,
				Size:   13,
				Color:  color.NewWhite(),
				Family: consts.Arial,
				Style:  consts.Bold,
				Align:  consts.Center,
			})
		})
	})

	m.SetBackgroundColor(color.NewWhite())

	m.Row(5, func() {
		m.Col(12, func() {
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("This medical certificate was issued on %s.", formatDate(certificateData.StartDate)), props.Text{
				Top:    2,
				Size:   11,
				Family: consts.Arial,
			})
		})
	})

	m.Row(15, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("We would like to inform that the patient was diagnosed with %s, according to the symptoms explained by the patient during the consultation with us.", certificateData.Diagnosis), props.Text{
				Top:             2,
				Size:            11,
				Family:          consts.Arial,
				VerticalPadding: 2,
			})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("Therefore, we recommend the patient to take a rest, starting from %s until %s.", formatDate(certificateData.StartDate), formatDate(certificateData.EndDate)), props.Text{
				Top:    2,
				Size:   11,
				Family: consts.Arial,
			})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("To whom it may concern, please take information above into consideration.", props.Text{
				Top:    2,
				Size:   11,
				Family: consts.Arial,
			})
		})
	})

	m.Row(5, func() {
		m.Col(12, func() {
		})
	})

	m.Row(5, func() {
		m.Col(12, func() {
			m.Row(10, func() {
				m.Col(12, func() {
					m.Text("Consulted with", props.Text{
						Top:    2,
						Size:   11,
						Family: consts.Arial,
						Align:  consts.Right,
					})
				})
			})

			m.Row(10, func() {
				m.Col(12, func() {
					m.Text(certificateData.DoctorName, props.Text{
						Top:    2,
						Size:   11,
						Family: consts.Arial,
						Align:  consts.Right,
						Style:  consts.Bold,
					})
				})
			})
		})
	})
}

func buildPrescription(m pdf.Maroto, prescriptionData entities.PrescriptionData) {
	headings := getHeadings()
	contents := [][]string{{"Name", fmt.Sprintf(" : %s", prescriptionData.PatientName)}, {"Gender", fmt.Sprintf(" : %s", prescriptionData.PatientGender.Name)}, {"Birth Date", fmt.Sprintf(" : %s", formatDate(prescriptionData.PatientBirthDate))}, {"Age", fmt.Sprintf(" : %d y.o.", prescriptionData.PatientAge)}}

	m.SetBackgroundColor(getTealColor())
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Patient Data", props.Text{
				Top:    2,
				Size:   13,
				Color:  color.NewWhite(),
				Family: consts.Arial,
				Style:  consts.Bold,
				Align:  consts.Center,
			})
		})
	})

	m.SetBackgroundColor(color.NewWhite())

	m.TableList(headings, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      10,
			GridSizes: []uint{3, 3},
		},
		ContentProp: props.TableListContent{
			Size:      10,
			GridSizes: []uint{3, 3},
		},
		Align:                  consts.Left,
		HeaderContentSpace:     1,
		Line:                   false,
		VerticalContentPadding: 2,
	})

	m.Row(10, func() {
		m.Col(12, func() {
		})
	})

	m.SetBackgroundColor(getTealColor())
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Prescription", props.Text{
				Top:    2,
				Size:   13,
				Color:  color.NewWhite(),
				Family: consts.Arial,
				Style:  consts.Bold,
				Align:  consts.Center,
			})
		})
	})

	m.SetBackgroundColor(color.NewWhite())

	m.Row(5, func() {
		m.Col(12, func() {
		})
	})

	contentsProducts := [][]string{}

	for i := 0; i < len(prescriptionData.Products); i++ {
		contentsProducts = append(contentsProducts, []string{prescriptionData.Products[i].Name, fmt.Sprintf("%d Per %s", prescriptionData.Quantities[i], prescriptionData.Products[i].SellingUnit)})
	}

	m.TableList([]string{"Product Name", "Quantity"}, contentsProducts, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      10,
			GridSizes: []uint{6, 6},
		},
		ContentProp: props.TableListContent{
			Size:      10,
			GridSizes: []uint{6, 6},
		},
		Align:                  consts.Left,
		HeaderContentSpace:     1,
		Line:                   false,
		VerticalContentPadding: 2,
	})

	m.Row(5, func() {
		m.Col(12, func() {
		})
	})

	m.Row(5, func() {
		m.Col(12, func() {
			m.Row(10, func() {
				m.Col(12, func() {
					m.Text("Consulted with", props.Text{
						Top:    2,
						Size:   11,
						Family: consts.Arial,
						Align:  consts.Right,
					})
				})
			})

			m.Row(10, func() {
				m.Col(12, func() {
					m.Text(prescriptionData.DoctorName, props.Text{
						Top:    2,
						Size:   11,
						Family: consts.Arial,
						Align:  consts.Right,
						Style:  consts.Bold,
					})
				})
			})
		})
	})
}

func getHeadings() []string {
	return []string{"", ""}
}

func getTealColor() color.Color {
	return color.Color{
		Red:   3,
		Green: 166,
		Blue:  166,
	}
}

func formatDate(timestamp string) string {
	timeArr := strings.Split(timestamp, "-")
	date, _ := strconv.Atoi(timeArr[2])
	return fmt.Sprintf("%s %d %s", constants.MonthsMap[timeArr[1]], date, timeArr[0])
}
